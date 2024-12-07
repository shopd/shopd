package router

import "github.com/a-h/templ"

type TemplRegistry struct {
	// api is a map of registered templ components for api routes
	api map[string]templ.Component
	// static is a map of registered templ components for static site routes
	static map[string]templ.Component
}

func (tr *TemplRegistry) RegisterAPI(route string, c templ.Component) {
	if _, registered := tr.api[route]; !registered {
		tr.api[route] = c
	}
}

func (tr *TemplRegistry) RegisterStatic(route string, c templ.Component) {
	if _, registered := tr.static[route]; !registered {
		tr.static[route] = c
	}
}

var RegisterAPI func(tr *TemplRegistry) = nil

var RegisterStatic func(tr *TemplRegistry) = nil

func NewTemplRegistry() (tr *TemplRegistry) {
	tr = &TemplRegistry{
		api:    make(map[string]templ.Component),
		static: make(map[string]templ.Component),
	}

	// Register funcs are optionally defined on init
	if RegisterAPI != nil {
		RegisterAPI(tr)
	}
	if RegisterStatic != nil {
		RegisterStatic(tr)
	}

	return tr
}
