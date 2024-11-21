package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
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

	// TODO Build all dev artifacts.
	// Do not call Clean here, otherwise cmd.SyncData has to run
	// each time the dev server is started, that would be slow
	// err = BuildDev()
	// if err != nil {
	// 	return err
	// }

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
	// err = devSite(txSession(EnvDev))
	// if err != nil {
	// 	return err
	// }

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
