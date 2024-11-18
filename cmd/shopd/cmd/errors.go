package cmd

import (
	"github.com/mozey/errors"
)

var ErrCmd = errors.NewCause("cmd")

var ErrNotImplemented = func(msg string) error {
	return errors.NewWithCausef(ErrCmd, "not implemented %s", msg)
}
