// Code generated with https://github.com/mozey/config DO NOT EDIT

package config

import (
	"bytes"
	"text/template"
)

// ExecTemplateDomainDir fills APP_TEMPLATE_DOMAIN_DIR with the given params
func (c *Config) ExecTemplateDomainDir() string {
	t := template.Must(template.New("templateDomainDir").Parse(c.templateDomainDir))
	b := bytes.Buffer{}
	_ = t.Execute(&b, map[string]interface{}{

		"Dir":    c.dir,
		"Domain": c.domain,
	})
	return b.String()
}
