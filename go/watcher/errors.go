package watcher

import (
	"github.com/mozey/errors"
)

var ErrWatcher = errors.NewCause("watcher")

var ErrAbsPath = func(p string) error {
	return errors.NewWithCausef(ErrWatcher, "path is not absolute %s", p)
}

var ErrRecursion = func(r int) error {
	return errors.NewWithCausef(ErrWatcher, "recursion limit exceeded %d", r)
}
