package main

import (
	"fmt"

	"github.com/mozey/errors"
)

var ErrMage = errors.NewCause("mage")

var ErrDep = func(dep string) error {
	return errors.NewWithCausef(ErrMage, "install dependency %s", dep)
}

var ErrDepCmd = func(cmd string) error {
	return errors.NewWithCausef(ErrMage, "invalid dependency %s", cmd)
}

var ErrParamInvalid = func(param, value string) error {
	msg := param
	if value != "" {
		msg = fmt.Sprintf("%s %s", param, value)
	}
	return errors.NewWithCausef(ErrMage, "%s is invalid", msg)
}