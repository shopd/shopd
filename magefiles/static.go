package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shopd/shopd/go/fileutil"
)

// mageStaticPath must not put static binary in APP_DIR,
// might cause confusion with global mage
func mageStaticPath() string {
	return filepath.Join(conf.Dir(), "magefiles", "mage")
}

func mageCmd(task string, args ...string) (cmd string, err error) {
	staticPath := mageStaticPath()
	if len(args) > 0 {
		return fmt.Sprintf("%s %s %s", staticPath, task, args[0]), nil
	}
	return fmt.Sprintf("%s %s", staticPath, task), nil
}

// BuildStaticMage builds a static executable for mage.
// Recompiling mage for each tmux pane will error,
// due to a race condition with go build?
// Static executable is also more responsible for dev
func BuildStaticMage() (err error) {
	staticPath := mageStaticPath()
	if fileutil.PathExists(staticPath) {
		// Remove static executable if it exist
		err = os.Remove(staticPath)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// Compile static executable
	cmd := exec.Command(mage, "-compile", staticPath)
	err = printCombinedOutput(cmd)
	if err != nil {
		return errors.WithStack(err)
	}

	// Make sure static executable exists
	if !fileutil.PathExists(staticPath) {
		return errors.WithStack(ErrMageStatic(staticPath))
	}

	return nil
}
