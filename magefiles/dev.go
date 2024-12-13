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

	// Caddy gateway
	err = devCaddy(txSession(EnvDev))
	if err != nil {
		return err
	}

	// Pane for shopd backend
	// err = tmuxSplitWindow(fmt.Sprintf("%s:0.0", txSession(EnvDev)), txVertical)
	err = tmuxSplitWindow(fmt.Sprintf("%s:0.0", txSession(EnvDev)), txHorizontal)
	if err != nil {
		return err
	}

	// ...........................................................................
	// Watch for...

	// ...site changes (api and content templates)
	err = tmuxNewWindow(txSession(EnvDev))
	if err != nil {
		return err
	}
	err = devSite(txSession(EnvDev))
	if err != nil {
		return err
	}

	// ...tailwind changes
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
	err = devTailwind(txSession(EnvDev))
	if err != nil {
		return err
	}

	// ...app changes
	err = tmuxSelectPane(fmt.Sprintf("%s:1.2", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = devApp(txSession(EnvDev))
	if err != nil {
		return err
	}

	// ...........................................................................
	// Watch for shopd changes
	err = tmuxSelectWindow(fmt.Sprintf("%s:0", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = tmuxSelectPane(fmt.Sprintf("%s:0.1", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = devShopd(txSession(EnvDev))
	if err != nil {
		return err
	}

	// ...........................................................................
	// Default view
	err = tmuxSelectWindow(fmt.Sprintf("%s:0", txSession(EnvDev)))
	if err != nil {
		return err
	}
	err = tmuxSelectPane(fmt.Sprintf("%s:0.1", txSession(EnvDev)))
	if err != nil {
		return err
	}

	return nil
}

func appWatcherCmd() string {
	appEntryPath := filepath.Join(conf.Dir(), "src", "app.ts")
	outDir := filepath.Join(conf.Dir(), "build")
	return fmt.Sprintf(`pnpx esbuild %s \
		--bundle --outdir=%s --watch`, appEntryPath, outDir)
}

// devApp runs the web app watcher in a tmux session
func devApp(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))

	log.Info().Msg("Web app live-reload watcher...")
	err = tmuxSendCmd(session, appWatcherCmd())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func tailwindWatcherCmd() string {
	appEntryPath := filepath.Join(conf.Dir(), "src", "app.css")
	outPath := filepath.Join(conf.Dir(), "build", "app.css")
	return fmt.Sprintf(`pnpx tailwindcss \
		-i %s -o %s \
		--minify --watch`, appEntryPath, outPath)
}

// devTailwind runs the tailwind watcher in a tmux session
func devTailwind(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))

	log.Info().Msg("Tailwind live-reload watcher...")
	err = tmuxSendCmd(session, tailwindWatcherCmd())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func shopdWatcherCmd() string {
	binPath := filepath.Join(conf.Dir(), "build", "bin")
	mainPath := filepath.Join(conf.Dir(), "cmd", "shopd", "main.go")
	return fmt.Sprintf(`%s \
	--build.cmd "go build -o %s %s" \
	--build.bin "%s run" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.exclude_dir "magefiles" \
	--build.exclude_dir "vendor" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true`, air, binPath, mainPath, binPath)
}

// devShopd runs shopd backend in a tmux session
func devShopd(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))
	mg.Deps(mg.F(Dep, air))

	log.Info().Msg("Dev shopd backend live-reload watcher...")
	err = tmuxSendCmd(session, shopdWatcherCmd())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// templWatcherCmd creates a command that uses the built-in templ watcher.
// It rebuilds when templ files are changed and then calls the specified cmd.
// Inspired by the Makefile example here
// https://templ.guide/commands-and-tools/live-reload-with-other-tools/#setting-up-the-makefile
// TODO Instead of the templ live reload proxy (doesn't work with https?),
// rather make the web app poll the server version for reloads
// https://templ.guide/commands-and-tools/live-reload
func templWatcherCmd() string {
	return fmt.Sprintf("%s generate -v --watch --path %s --cmd \"%s\"",
		templ,
		filepath.Join(conf.Dir(), "www"),
		mageCmd("BuildSite", "dev"))
}

// devSite runs a watcher for the templ files in a tmux session.
// The --watch flag makes templ generate *_templ.txt files.
// The HTML is read from the txt files as a dev optimisation,
// instead of being embedded in the *_templ.go files
// https://github.com/a-h/templ/pull/366
func devSite(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))
	mg.Deps(mg.F(Dep, templ))

	log.Info().Msg("Templ live-reload watcher...")

	err = tmuxSendCmd(session, templWatcherCmd())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// devCaddy runs caddy in a tmux session
func devCaddy(session string) (err error) {
	mg.Deps(mg.F(Dep, tmux))

	log.Info().Msg("Caddy gateway...")

	err = CaddyfileGenDev()
	if err != nil {
		return errors.WithStack(err)
	}

	caddyCmd := fmt.Sprintf(
		"caddy run --config %s", filepath.Join(conf.Dir(), "Caddyfile"))
	err = tmuxSendCmd(session, caddyCmd)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
