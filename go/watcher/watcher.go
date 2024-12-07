package watcher

import (
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/shopd/shopd-proto/go/fileutil"
)

// recursionLimit when watching directories recursively
const recursionLimit = 10

// defaultDelayMS is the delay in milliseconds before the change callback
const defaultDelayMS = 1500

type ChangeCB func(path string)

type Watcher struct {
	change         ChangeCB
	delayMS        int
	shutdown       chan os.Signal
	watcher        *fsnotify.Watcher
	includePaths   []string
	excludePaths   []*regexp.Regexp
	includeChanges []*regexp.Regexp
	excludeChanges []*regexp.Regexp
}

type WatcherParams struct {
	// Change callback is called when a watched file is modified
	Change ChangeCB
	// DelayMS is the delay in milliseconds before the change callback is called
	DelayMS int
	// IncludePaths lists paths to include recursively,
	// changes to files in included directories will trigger change events
	IncludePaths []string
	// ExcludePaths lists patterns to ignore when including directories to monitor
	ExcludePaths []string
	// IncludeChanges lists path patterns to match on modification events
	IncludeChanges []string
	// ExcludeChanges lists path patterns to ignore on modification events
	ExcludeChanges []string
}

func NewWatcher(params WatcherParams) (h *Watcher, err error) {
	h = &Watcher{}

	h.change = params.Change

	delayMS := defaultDelayMS
	if params.DelayMS > 0 {
		// Override default
		delayMS = params.DelayMS
	}
	h.delayMS = delayMS

	h.includePaths = params.IncludePaths
	h.excludePaths = make([]*regexp.Regexp, 0)
	h.includeChanges = make([]*regexp.Regexp, 0)
	h.excludeChanges = make([]*regexp.Regexp, 0)

	// Path patterns to exclude
	for _, pattern := range params.ExcludePaths {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return h, errors.WithStack(err)
		}
		h.excludePaths = append(h.excludePaths, r)
	}

	// Change patterns to include
	for _, pattern := range params.IncludeChanges {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return h, errors.WithStack(err)
		}
		h.includeChanges = append(h.includeChanges, r)
	}

	// Change patterns to exclude
	for _, pattern := range params.ExcludeChanges {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return h, errors.WithStack(err)
		}
		h.excludeChanges = append(h.excludePaths, r)
	}

	return h, nil
}

// excludePath returns true if the path must be excluded
func (h *Watcher) excludePath(p string) bool {
	for _, r := range h.excludePaths {
		if r.Match([]byte(p)) {
			return true
		}
	}
	return false
}

// includeChange returns true if the change path must be included
func (h *Watcher) includeChange(p string) bool {
	if strings.TrimSpace(p) == "" {
		// Path is sometimes empty on change event
		return true
	}
	for _, r := range h.includeChanges {
		if r.Match([]byte(p)) {
			return true
		}
	}
	return false
}

// excludeChange returns true if the change path must be excluded
func (h *Watcher) excludeChange(p string) bool {
	if strings.TrimSpace(p) == "" {
		// Path is sometimes empty on change event
		return true
	}
	for _, r := range h.excludeChanges {
		if r.Match([]byte(p)) {
			return true
		}
	}
	return false
}

// setup fsnotify watcher as per config on the handler
func (h *Watcher) setup() (err error) {
	// Include paths
	for _, p := range h.includePaths {
		if !filepath.IsAbs(p) {
			return errors.WithStack(ErrAbsPath(p))
		}
		r := 0
		err = filepath.Walk(p,
			func(p string, info os.FileInfo, err error) error {
				if r > recursionLimit {
					errors.WithStack(ErrRecursion(r))
				}
				if !h.excludePath(p) {
					// Only watch directories
					if fileutil.IsDir(p) {
						// fmt.Println("add", p)
						h.watcher.Add(p)
					}
				}
				r++
				return nil
			})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// Timeout executes the callback after the specified duration.
// Abort execution by sending a message on the cancel chan
func Timeout(cancel chan bool, duration time.Duration, cb func()) {
	if cancel == nil {
		return
	}
	select {
	case <-time.After(duration):
		cb()
	case <-cancel:
		// Abort
	}
}

// Run the watcher, ctrl + c to exit
func (h *Watcher) Run() (err error) {
	// Relay incoming signal from user (ctrl + c) to channel h.sig
	// See comments in cmd/shopd/cmd/run.go re. graceful shutdown
	h.shutdown = make(chan os.Signal, 1)
	signal.Notify(h.shutdown, os.Interrupt, syscall.SIGTERM)

	// Create new watcher
	h.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return errors.WithStack(err)
	}
	defer h.watcher.Close()
	err = h.setup()
	if err != nil {
		return err
	}

	// Start listening for events
	// https://github.com/fsnotify/fsnotify?tab=readme-ov-file#usage
	var watcherErr error = nil
	go func() {
		var cancel chan bool
		for {
			select {
			case event, more := <-h.watcher.Events:
				if !more {
					// Channel closed and all values already received
					// https://gobyexample.com/closing-channels
					h.Signal()
					return
				}
				// fmt.Printf("event %#v\n", event)
				if event.Has(fsnotify.Write) {
					match := true
					if len(h.includeChanges) > 0 {
						// If any inclusion patterns are set,
						// then the path must match one of them
						match = h.includeChange(event.Name)
					}
					// Path must not be excluded
					match = match && !h.excludeChange(event.Name)
					if match {
						// Reset cancel chan
						cancel = make(chan bool)
						// Use a timeout in case multiple files were changed
						go Timeout(cancel, time.Duration(h.delayMS)*time.Millisecond,
							func() {
								h.change(event.Name)
							})
					}
				}
			case err, more := <-h.watcher.Errors:
				if !more {
					h.Signal()
					return
				}
				// Return on first error
				watcherErr = err
				h.Signal()
				return
			}
		}
	}()

	<-h.shutdown
	return watcherErr
}

// Signal on the shutdown channel
func (h *Watcher) Signal() {
	// "SIGINT is usually user-initiated,
	// while SIGTERM can be system or process-initiated"
	h.shutdown <- os.Signal(syscall.SIGTERM)
}
