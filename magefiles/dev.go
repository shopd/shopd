package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/fileutil"
)

// Dev starts the dev server.
// To avoid generating static files with test data,
// the "sync data" cmd must be executed manually,
// make use of the DomainGen target
func Dev() (err error) {
	return dev()
}

func dev() (err error) {
	mg.Deps(mg.F(Dep, tmux))

	// Err if session exists
	if tmuxSessionExists(txSession(EnvDev)) {
		return errors.WithStack(ErrSessionExists(txSession(EnvDev)))
	}

	// Build static binary for mage
	err = BuildStaticMage()
	if err != nil {
		return err
	}

	// Remove dev emails
	emailDir := filepath.Join(conf.Dir(), "email")
	err = os.RemoveAll(emailDir)
	if err != nil {
		return errors.WithStack(err)
	}
	err = fileutil.MkdirAll(emailDir)
	if err != nil {
		return errors.WithStack(err)
	}

	// Build all dev artefacts.
	// Do not call Clean here, otherwise cmd.SyncData has to run
	// each time the dev server is started, that would be slow
	err = BuildDev()
	if err != nil {
		return err
	}

	// ...........................................................................
	// Create tmux session
	err = tmuxNewSession(txSession(EnvDev))
	if err != nil {
		return err
	}

	// TODO Caddy gateway
	// err = devCaddy(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }

	// Pane for shopd backend
	err = tmuxSplitWindow(txSession(EnvDev), txVertical)
	if err != nil {
		return err
	}
	// Pane for shopd backend watcher
	err = tmuxSplitWindow(fmt.Sprintf("%s:0.0", txSession(EnvDev)), txHorizontal)
	if err != nil {
		return err
	}

	// ...........................................................................
	// Watch for...

	// TODO ...site changes
	err = tmuxNewWindow(txSession(EnvDev))
	if err != nil {
		return err
	}
	// TODO Remove this?
	// err = devSite(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }
	err = devTempl(txSession(EnvDev))
	if err != nil {
		return err
	}

	// TODO ...tailwind changes
	err = tmuxSplitWindow(txSession(EnvDev), txVertical)
	if err != nil {
		return err
	}
	err = tmuxSelectPane(fmt.Sprintf("%s:1.0", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = tmuxSplitWindow(txSession(EnvDev), txHorizontal)
	if err != nil {
		return err
	}
	// err = devTailwind(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }

	// TODO ...app changes
	err = tmuxSelectPane(fmt.Sprintf("%s:1.2", txSession(EnvDev)))
	if err != nil {
		return err
	}
	// err = devApp(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }

	// ...........................................................................
	// TODO Watch for shopd changes
	err = tmuxSelectWindow(fmt.Sprintf("%s:0", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = tmuxSelectPane(fmt.Sprintf("%s:0.1", txSession(EnvDev)))
	if err != nil {
		return err
	}
	// err = devShopd(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }

	// ...........................................................................
	// Default view
	err = tmuxSelectWindow(fmt.Sprintf("%s:0", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = tmuxSelectPane(fmt.Sprintf("%s:0.2", txSession(EnvDev)))
	if err != nil {
		return err
	}

	return nil
}

// TODO Remove this?
// // devSite runs WatchSite in a tmux session
// func devSite(session string) (err error) {
// 	mg.Deps(mg.F(Dep, tmux))

// 	log.Info().Msg("Static site helper live-reload watcher...")
// 	cmd, err := mageCmd(taskWatchSite)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	err = tmuxSendCmd(session, cmd)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	return nil
// }

func templStaticGenCmd() string {
	staticGenCmd := fmt.Sprintf("go run %s/... static gen --env dev",
		filepath.Join(conf.Dir(), "cmd", "shopd"))

	cmd := fmt.Sprintf("%s generate -v --watch --path %s --cmd \"%s\"",
		templ,
		filepath.Join(conf.Dir(), "www"),
		staticGenCmd)

	return cmd
}

// TODO Remove this, see README.md
func DebugTemplStaticGen() {
	fmt.Println(templStaticGenCmd())
}

// devTempl runs a watcher for the templ files in a tmux session.
// The --watch flag makes templ generate *_templ.txt files.
// The HTML is read from the txt files as a dev optimisation,
// instead of being embedded in the *_templ.go files
// https://github.com/a-h/templ/pull/366
func devTempl(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))
	mg.Deps(mg.F(Dep, templ))

	log.Info().Msg("Templ live-reload watcher...")

	err = tmuxSendCmd(session, templStaticGenCmd())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// TODO Remove this?
// // WatchSite rebuilds and static site helper
// func WatchSite() (err error) {
// 	w, err := watcher.NewWatcher(watcher.WatcherParams{
// 		Change: func(p string) {
// 			envMap, err := versionEnv()
// 			if err != nil {
// 				log.Error().Stack().Err(err).Msg("")
// 				return
// 			}

// 			log.Info().Str("path", p).Msg("WatchSite")

// 			err = buildSite(EnvDev, envMap, false)
// 			if err != nil {
// 				log.Error().Stack().Err(err).Msg("")
// 				return
// 			}
// 		},
// 		DelayMS: 0, // Using default delay
// 		IncludePaths: []string{
// 			filepath.Join(conf.Dir(), "www", "components"),
// 			filepath.Join(conf.Dir(), "www", "content"),
// 			// TODO Domain can override?
// 			// filepath.Join(conf.ExecTemplateDomainDir(), "www", "content"),
// 		},
// 		ExcludePaths:   []string{".*public.*"},
// 		ExcludeChanges: []string{".*hugo_build.lock$"},
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return w.Run()
// }
