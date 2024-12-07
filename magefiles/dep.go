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
	air     = "air"
	caddy   = "caddy"
	configu = "configu"
	find    = "find"
	git     = "git"
	goCmd   = "go"
	mage    = "mage"
	sqlc    = "sqlc"
	templ   = "templ"
	tmux    = "tmux"
)

// Dep checks for programs used with mage
func Dep(cmd string) error {
	switch cmd {
	case air:
		err := exec.Command("air", "-v").Run()
		if err != nil {
			// cd ~ && go install github.com/air-verse/air@v1.61.1
			return errors.WithStack(ErrDep("https://github.com/air-verse/air"))
		}
	case caddy:
		err := exec.Command("caddy", "version").Run()
		if err != nil {
			// Caddy also requires Network Security Services
			// brew install nss
			return errors.WithStack(ErrDep("https://formulae.brew.sh/formula/caddy"))
		}
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
	case git:
		err := exec.Command("git", "version").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://git-scm.com/"))
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
	case templ:
		err := exec.Command("templ", "version").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://templ.guide/quick-start/installation"))
		}
	case tmux:
		err := exec.Command("tmux", "-V").Run()
		if err != nil {
			return errors.WithStack(ErrDep("https://formulae.brew.sh/formula/tmux"))
		}
	default:
		return errors.WithStack(ErrDepCmd(cmd))
	}
	return nil
}
