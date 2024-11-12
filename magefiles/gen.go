package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/shopd/shopd/go/fileutil"
	"github.com/shopd/shopd/go/share"
)

var domainConfigGen = "example.com"

// ConfigGen generates the config helper package
func ConfigGen() (err error) {
	mg.Deps(mg.F(Dep, configu))
	mg.Deps(mg.F(Dep, goCmd))

	// Requires a config file to generate the helpers
	err = EnvGen("dev", domainConfigGen)
	if err != nil {
		return err
	}

	// Generate helper package
	cmd := exec.Command(configu, "-generate", filepath.Join("go", "config"))
	err = printCombinedOutput(cmd)
	if err != nil {
		return err
	}

	// Format code
	cmd = exec.Command("go", "fmt",
		filepath.Join(conf.Dir(), "go", "config", "config.go"))
	printCombinedOutput(cmd) // Ignore errors

	cmd = exec.Command("go", "fmt",
		filepath.Join(conf.Dir(), "go", "config", "fn.go"))
	printCombinedOutput(cmd) // Ignore errors

	cmd = exec.Command("go", "fmt",
		filepath.Join(conf.Dir(), "go", "config", "template.go"))
	printCombinedOutput(cmd) // Ignore errors

	return nil
}

func escapeDomain(env, domain string) string {
	// Prefix env and replace dots with hyphens, e.g. dev-example-com
	return fmt.Sprintf("%s-%s", env, share.EscapeDomain(domain))
}

// EnvGen generates .env files for shopd from templates.
// Do not make use of conf global in this target,
// until after the config file is generated
func EnvGen(env, domain string) (err error) {
	appDir := os.Getenv(paramAppDir)

	escapedDomainName := escapeDomain(env, domain)
	envFile := filepath.Join(appDir, fmt.Sprintf(".env.%s.sh", escapedDomainName))
	fmt.Println("envFile", envFile)

	buf := bytes.Buffer{}
	buf.WriteString("{}")
	err = fileutil.WriteBytes(envFile, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
