package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	OSDarwin  = "darwin"
	OSLinux   = "linux"
	OSWindows = "windows"
	OSOther   = "other"
)

// detectOS detects the host OS.
// Return values must correspond to GOOS listed here
// https://go.dev/doc/install/source#environment
func detectOS() (os string, err error) {
	cmd := exec.Command("uname", "-s")
	b, err := cmd.Output()
	if err != nil {
		return os, errors.WithStack(err)
	}
	s := string(b)

	if strings.Contains(s, "Darwin") {
		return OSDarwin, nil
	}
	if strings.Contains(s, "Linux") {
		return OSLinux, nil
	}
	if strings.Contains(s, "CYGWIN") ||
		strings.Contains(s, "MINGW32") ||
		strings.Contains(s, "MSYS") ||
		strings.Contains(s, "MINGW") {
		return OSWindows, nil
	}

	log.Info().Str("output", s).Msg("")
	return OSOther, nil
}

// DetectOS prints the GOOS value for the host OS
func DetectOS() (err error) {
	os, err := detectOS()
	if err != nil {
		return err
	}
	fmt.Println(os)
	return nil
}

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
