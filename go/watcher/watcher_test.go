package watcher_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd-proto/go/fileutil"
	"github.com/shopd/shopd-proto/go/testutil"
	"github.com/shopd/shopd-proto/go/watcher"
)

func TestWatcher_Run(t *testing.T) {
	is := testutil.Setup(t)

	// Setup tmp dir and files
	tmp, err := os.MkdirTemp("", t.Name())
	is.NoErr(err)
	defer func() {
		// TODO Remove tmp dir
		log.Info().Str("path", tmp).Msg("tmp dir")
		// os.RemoveAll(tmp)
	}()

	sub1 := filepath.Join(tmp, "sub1")
	err = os.Mkdir(sub1, fileutil.PermDirDefault)
	is.NoErr(err)
	sub2 := filepath.Join(sub1, "sub2")
	err = os.Mkdir(sub2, fileutil.PermDirDefault)
	is.NoErr(err)
	sub3 := filepath.Join(sub1, "sub3")
	err = os.Mkdir(sub3, fileutil.PermDirDefault)
	is.NoErr(err)

	file1 := filepath.Join(sub1, "file1")
	err = fileutil.WriteBytes(file1, []byte("init"))
	is.NoErr(err)
	file2 := filepath.Join(sub2, "file2")
	err = fileutil.WriteBytes(file2, []byte("init"))
	is.NoErr(err)
	file3 := filepath.Join(sub3, "file3")
	err = fileutil.WriteBytes(file3, []byte("init"))
	is.NoErr(err)

	// Run watcher
	delay1 := 1
	delay10 := 10
	changeMap := make(map[string]bool)
	var wg sync.WaitGroup
	w, err := watcher.NewWatcher(watcher.WatcherParams{
		Change: func(p string) {
			log.Info().Str("path", p).Msg("change")
			changeMap[p] = true
		},
		DelayMS:        delay1,
		IncludePaths:   []string{tmp},
		ExcludePaths:   []string{".*sub3.*"},
		IncludeChanges: []string{".*file.*"},
		ExcludeChanges: []string{".*file2.*"},
	})
	is.NoErr(err)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = w.Run()
		is.NoErr(err)
		log.Info().Msg("watcher exited")
	}()

	// Delay before making changes,
	// to give watcher time to initialise
	time.Sleep(time.Duration(delay1) * time.Millisecond)

	// Make some changes
	err = fileutil.WriteBytes(file1, []byte("change"))
	is.NoErr(err)
	err = fileutil.WriteBytes(file2, []byte("change"))
	is.NoErr(err)
	err = fileutil.WriteBytes(file3, []byte("change"))
	is.NoErr(err)

	time.Sleep(time.Duration(delay10) * time.Millisecond)
	// Signal shutdown
	w.Signal()
	// Wait for watcher to shutdown
	wg.Wait()

	// Confirm changes
	_, found := changeMap[file1]
	is.True(found)
	_, found = changeMap[file2]
	is.True(!found)
	_, found = changeMap[file3]
	is.True(!found)
}
