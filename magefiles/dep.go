package main

import (
	"os/exec"

	"github.com/pkg/errors"
)

const (
	configu = "configu"
	find    = "find"
	goCmd   = "go"
	sqlc    = "sqlc"
)

// Dep checks for programs used with mage
func Dep(cmd string) error {
	switch cmd {
	case configu:
		err := exec.Command("configu", "--help").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://github.com/mozey/config"))
		}
	case find:
		err := exec.Command("find", ".", "-quit").Run()
		if err != nil {
			return errors.WithStack(ErrDep("find"))
		}
	case goCmd:
		err := exec.Command("go", "version").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://golang.org"))
		}
	case sqlc:
		err := exec.Command("sqlc", "version").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://formulae.brew.sh/formula/sqlc"))
		}
	default:
		return errors.WithStack(ErrDepCmd(cmd))
	}
	return nil
}
