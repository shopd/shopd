package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/share"
)

const enter = "ENTER"

const txKillSession = "kill-session"

const txNewSession = "new-session"

const txNewWindow = "new-window"

const txSelectPane = "select-pane"

const txSelectWindow = "select-window"

const txSendKeys = "send-keys"

const txSplitWindow = "split-window"

func txSession(env string) string {
	// Dev might have multiple tmux sessions and
	// Caddy processes each with a Caddyfile
	return fmt.Sprintf("shopd-%s", escapeDomain(env, conf.Domain()))
}

func txSessionCaddy() string {
	return fmt.Sprintf("shopd-caddy")
}

type tmuxKeys string

const txCtrlC tmuxKeys = "C-c"

func targetShopd(env string) string {
	return fmt.Sprintf("%s:0.2", txSession(env))
}

const taskWatchApp = "watchApp"

const taskWatchTailwind = "watchTailwind"

const taskWatchShopd = "watchShopd"

const taskWatchSite = "watchSite"

// type tmuxKeys string

// const txCtrlC tmuxKeys = "C-c"

// tmuxSessionExists returns true if the session exists
func tmuxSessionExists(session string) bool {
	out, err := exec.Command(tmux, "ls").Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), session)
}

// tmuxNewSession creates a new tmux session
func tmuxNewSession(session string) (err error) {
	cmd := exec.Command(tmux, txNewSession, "-d", "-s", session)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxNewWindow creates a new tmux window
func tmuxNewWindow(session string) (err error) {
	cmd := exec.Command(tmux, txNewWindow, "-t", session)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxKillSession creates a new tmux window
func tmuxKillSession(session string) (err error) {
	cmd := exec.Command(tmux, txKillSession, "-t", session)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxSelectWindow creates a new tmux window
func tmuxSelectWindow(target string) (err error) {
	cmd := exec.Command(tmux, txSelectWindow, "-t", target)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxSelectPane creates a new tmux window
func tmuxSelectPane(target string) (err error) {
	cmd := exec.Command(tmux, txSelectPane, "-t", target)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxSendCmd sends a single command to the tmux target.
// See link re. target syntax: ${SESSION}:${WINDOW}.${PANE}
// https://superuser.com/a/492549/537059
func tmuxSendCmd(target string, tmuxCmd string) (err error) {
	cmd := exec.Command(tmux, txSendKeys, "-t", target, tmuxCmd, enter)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// tmuxSendKeys sends keys to the tmux target
func tmuxSendKeys(target string, keys tmuxKeys) (err error) {
	mg.Deps(mg.F(Dep, tmux))

	// Unlike tmuxSendCmd, the command below doesn't end with ENTER
	err = exec.Command(
		tmux, txSendKeys, "-t", target, mage, string(keys)).Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

const txVertical = "-v"

const txHorizontal = "-h"

// tmuxSplitWindow splits the tmux window
func tmuxSplitWindow(target string, orientation string) (
	err error) {

	// Just split 50% to avoid inconsistencies between tmux versions
	cmd := exec.Command(
		tmux, txSplitWindow, orientation,
		"-t", target)
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	return nil
}

// Down stops server processes for the specified env
func Down(env string) (err error) {
	mg.Deps(mg.F(Dep, caddy))
	mg.Deps(mg.F(Dep, tmux))

	// Don't stop Caddy on prod server,
	// other domains might be hosted on the same server
	if env == share.EnvProd {
		// TODO Delete caddy config for domain
		// err = CaddyDelete()
		// if err != nil {
		// 	log.Error().Stack().Err(err).Msg("")
		// }

	} else {
		// Stop Caddy
		cmd := exec.Command(caddy, "stop")
		err = printCombinedOutput(cmd)
		if err == nil {
			log.Error().Stack().Err(err).Msg("")
		}
	}

	// Stop domain backend service
	err = tmuxKillSession(txSession(env))
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}

	return nil
}
