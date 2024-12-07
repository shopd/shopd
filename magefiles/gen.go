package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/fileutil"
	"github.com/shopd/shopd/go/share"
)

var domainConfigGen = "example.com"

// ConfigGen generates the config helper package
func ConfigGen() (err error) {
	mg.Deps(mg.F(Dep, configu))
	mg.Deps(mg.F(Dep, goCmd))

	// Requires a config file to generate the helpers
	env := share.EnvDev
	domain := domainConfigGen
	err = EnvGen(env, domain)
	if err != nil {
		return err
	}

	// Generate helper package
	escapedDomainName := escapeDomain(env, domain)
	cmd := exec.Command(configu,
		"-env", escapedDomainName,
		"-generate", filepath.Join("go", "config"))
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
	var profile, listen, portCaddy string
	sample := envTemplate
	if env == share.EnvDev {
		profile = "aws-local"
		listen = "https://localhost"
		portCaddy = ":8443" // TODO Find unused port

	} else if env == share.EnvProd {
		profile = "shopd"
		listen = fmt.Sprintf("https://%s", domain)
		portCaddy = ":443"

	} else {
		return errors.WithStack(ErrParamInvalid(share.ParamEnv, env))
	}

	escapedDomainName := escapeDomain(env, domain)
	envFile := filepath.Join(appDir, fmt.Sprintf(".env.%s.sh", escapedDomainName))
	fmt.Println("envFile", envFile)
	portAPI := ":8500" // TODO Find unused port

	// TODO Review this
	// Assuming macOS is used for development only,
	// and override prod settings accordingly.
	// This facilitates testing of the prod targets
	os, err := detectOS()
	if err != nil {
		return err
	}
	if os == OSDarwin {
		listen = "https://localhost"
		portCaddy = ":8443"
	}

	// TODO Parse domain settings file.
	// Can't make use of ExecTemplateDomainDir here
	settings := map[string]any{
		"Domain": domainConfigGen,
		"Hosts":  "",
	}

	// Templates snippets that require pre-rendering
	buf := bytes.Buffer{}
	t, err := template.New("Snip").Parse(envPortsTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	err = t.Execute(&buf, map[string]any{
		"API":   portAPI,
		"Caddy": portCaddy,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	portsPre := buf.String()

	// Execute template
	t, err = template.New("Env").Parse(sample)
	if err != nil {
		return errors.WithStack(err)
	}
	buf = bytes.Buffer{}
	err = t.Execute(&buf, map[string]any{
		"Listen":   listen,
		"Profile":  profile,
		"Settings": settings,
		"Templates": map[string]any{
			"AwsCreds":        "# AwsCreds",
			"ConfigTemplates": envConfigTemplatesTemplate,
			"Deps":            "# Deps",
			"Email":           "# Email",
			"Ext":             "# Ext",
			"Http":            "# Http",
			"Limits":          "# Limits",
			"Nats":            "# Nats",
			"Ports":           portsPre,
			"Session":         "# Session",
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	// Write generated file
	err = fileutil.WriteBytes(envFile, buf.Bytes())
	if err != nil {
		return err
	}
	log.Info().Str("file", envFile).Msg("Generated")

	return nil
}

const envPortsTemplate = `# Caddy is used as an API Gateway,
# and proxies requests to /api to this port
APP_PORT_API="{{.API}}"
APP_PORT_CADDY="{{.Caddy}}"`

const envConfigTemplatesTemplate = `# Config templates
APP_TEMPLATE_DOMAIN_DIR="{{.Dir}}/domains/{{.Domain}}"`

const envTemplate = `# Code generated with https://github.com/shopd/shopd

{{.Templates.AwsCreds}}

{{.Templates.Deps}}

APP_DOMAIN="{{.Settings.Domain}}"
APP_DOMAIN_HOSTS="{{.Settings.Hosts}}"

{{.Templates.Email}}

{{.Templates.Http}}

# InstanceID must be unique, like Domain, but it can be a shorter code
APP_INSTANCE_ID="shopd"

{{.Templates.Limits}}

# Connections as per config
APP_LISTEN="{{.Listen}}"

{{.Templates.Nats}}

{{.Templates.Ports}}

{{.Templates.Session}}

{{.Templates.ConfigTemplates}}

{{.Templates.Ext}}

AWS_PROFILE={{.Profile}}`
