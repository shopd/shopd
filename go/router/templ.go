package router

import "github.com/a-h/templ"

type TemplRegistry struct {
	// api is a map of registered templ components for api routes and methods
	api map[string]map[string]templ.Component
	// static is a map of registered templ components for static site routes
	static map[string]templ.Component
}

func (tr *TemplRegistry) RegisterAPI(route, method string, c templ.Component) {
	if _, registered := tr.api[route]; !registered {
		if _, registered := tr.api[route][method]; !registered {
			// Only register on the first call to this method
			tr.api[route][method] = c
		}
	}
}

func (tr *TemplRegistry) RegisterStatic(route string, c templ.Component) {
	if _, registered := tr.static[route]; !registered {
		// Only register on the first call to this method
		tr.static[route] = c
	}
}

func (tr *TemplRegistry) API(route, method string) (
	c templ.Component, err error) {

	methods, registered := tr.api[route]
	if !registered {
		return c, ErrRouteNotFound(route)
	}

	c, registered = methods[method]
	if !registered {
		return c, ErrRouteNotFound(route)
	}

	return c, nil
}

func (tr *TemplRegistry) Static(route string) (c templ.Component, err error) {
	c, registered := tr.static[route]
	if !registered {
		return c, ErrRouteNotFound(route)
	}
	return c, nil
}

var RegisterAPI func(tr *TemplRegistry) = nil

var RegisterStatic func(tr *TemplRegistry) = nil

func NewTemplRegistry() (tr *TemplRegistry) {
	tr = &TemplRegistry{
		api:    make(map[string]map[string]templ.Component),
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
