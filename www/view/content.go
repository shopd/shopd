package view

import "github.com/shopd/shopd-proto/go/share"

type Content struct {
	BaseURL      string
	DomainConfig share.DomainConfigExport
}
