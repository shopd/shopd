package view

import (
	"net/url"

	"github.com/shopd/shopd-proto/go/share"
)

type Content struct {
	baseURL      string
	DomainConfig share.DomainConfigExport
}

// BaseURL appends a path to baseURL
func (c *Content) BaseURL(p string) string {
	u, _ := url.Parse(c.baseURL)
	u.Path = p
	return u.String()
}

func NewContent() *Content {
	return &Content{
		baseURL: "https://localhost:8443/",
	}
}
