package view

import (
	"net/url"
	"path"

	"github.com/shopd/shopd-proto/go/share"
)

type Content struct {
	baseURL      string
	domainConfig share.DomainConfigExport
}

// BaseURL appends a path to baseURL
func (c *Content) BaseURL(p string) string {
	u, _ := url.Parse(c.baseURL)
	u.Path = p
	return u.String()
}

// StaticURL appends path to the static URL
func (c *Content) StaticURL(p string) string {
	u, _ := url.Parse(c.baseURL)
	u.Path = path.Join("s", p)
	return u.String()
}

func (c *Content) DomainConfig() share.DomainConfigExport {
	return c.domainConfig
}

type ContentParams struct {
	BaseURL      string
	DomainConfig share.DomainConfigExport
}

func NewContent(params ContentParams) Content {
	return Content{
		baseURL:      params.BaseURL,
		domainConfig: params.DomainConfig,
	}
}
