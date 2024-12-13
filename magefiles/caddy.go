package main

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/fileutil"
)

// CaddyfileGenDev generates the dev Caddyfile
func CaddyfileGenDev() (err error) {
	caddyfile := filepath.Join(conf.Dir(), "Caddyfile")

	// Sample template
	buf := bytes.Buffer{}
	t, err := template.New("Caddyfile").Parse(caddyFileTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	err = t.Execute(&buf, map[string]any{
		"Listen":    conf.Listen(),
		"PortCaddy": conf.PortCaddy(),
		"PortApi":   conf.PortApi(),
		"DomainDir": conf.ExecTemplateDomainDir(),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	// Caddyfile
	err = fileutil.WriteBytes(caddyfile, buf.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().Str("Caddyfile", caddyfile).Msg("Generated")

	return nil
}

const caddyFileTemplate = `# Code generated with https://github.com/shopd/shopd
{{.Listen}}{{.PortCaddy}} {
	# Using Caddy as an API Gateway

	# Routes on this path prefix is part of the standard API
	reverse_proxy /api* localhost{{.PortApi}}

	# Dynamically render static site content on dev
	reverse_proxy /* localhost{{.PortApi}}
}
`
