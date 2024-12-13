package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/cmd/shopd/cmd"
	"github.com/shopd/shopd/go/share"
)

// Do not hardcode env in code outside magefiles.
// Use go/config package to read env vars
const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

// Env vars starting with SHOPD_* are passed to sub-commands inline,
// these are not set in the parent shell or the config files
const (
	ShopdVersionBuild = "SHOPD_VERSION_BUILD"
	ShopdVersionSite  = "SHOPD_VERSION_SITE"
)

func appTargetDir() string {
	domainDir := conf.ExecTemplateDomainDir()
	return filepath.Join(domainDir, "www", "public", "app")
}

// BuildApp builds the TypeScript app
func BuildApp(env string) (err error) {
	return buildApp(env, appTargetDir())
}

func buildApp(env, targetDir string) (err error) {
	log.Info().Str("env", env).Msg("Building app")

	// TODO

	return nil
}

// BuildTailwind builds the css for the app
func BuildTailwind(env string) (err error) {
	return buildTailwind(env, appTargetDir())
}

func buildTailwind(env, targetDir string) (err error) {
	log.Info().Str("env", env).Msg("Building tailwind")

	// TODO

	return nil
}

// BuildSite builds the static site
func BuildSite(
	ctx context.Context, env string) (err error) {

	envMap, err := versionEnv()
	if err != nil {
		return err
	}

	err = buildSite(env, envMap)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func buildSite(
	env string, envMap []string) (err error) {

	log.Info().Str("env", env).Msg("Generate router init file for api")
	err = TemplApiGen(env)
	if err != nil {
		return err
	}

	if env == EnvDev {
		log.Info().Str("env", env).Msg("Generate router init file for site content")
		err = TemplSiteGen(env)
		if err != nil {
			return err
		}
	} else {
		return ErrNotImplemented(fmt.Sprintf("for env %s", env))
	}

	return nil
}

// gitRev return the current git revision
func gitRev() (rev string, err error) {
	cmd := exec.Command(
		git, "-C", conf.Dir(), "rev-parse", "--short", "--verify", "HEAD")
	out, err := combinedOutput(cmd)
	if err != nil {
		fmt.Println(string(out))
		return rev, err
	}
	return strings.TrimSpace(string(out)), nil
}

func versionInfo() (versionBuild, versionSite string, err error) {
	gitRev, err := gitRev()
	if err != nil {
		return versionBuild, versionSite, err
	}
	versionBuild = fmt.Sprintf("%s-%s", cmd.Version(), gitRev)
	versionSite = share.NowVersion()

	return versionBuild, versionSite, nil
}

func versionEnv() (envMap []string, err error) {
	versionBuild, versionSite, err := versionInfo()
	if err != nil {
		return envMap, errors.WithStack(err)
	}

	envMap = os.Environ()
	envMap = append(envMap,
		fmt.Sprintf("%s=%s", ShopdVersionBuild, versionBuild))
	envMap = append(envMap,
		fmt.Sprintf("%s=%s", ShopdVersionSite, versionSite))

	return envMap, nil
}

// BuildDev builds all the dev artefacts
func BuildDev() (err error) {
	envMap, err := versionEnv()
	if err != nil {
		return err
	}

	// Static site
	err = buildSite(EnvDev, envMap)
	if err != nil {
		return errors.WithStack(err)
	}

	// App artefacts
	err = BuildApp(EnvDev)
	if err != nil {
		return errors.WithStack(err)
	}
	err = BuildTailwind(EnvDev)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
