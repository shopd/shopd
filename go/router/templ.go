package router

import "github.com/a-h/templ"

type TemplRegistry struct {
	// api is a map of registered templ components for api routes and methods
	api map[string]map[string]templ.Component
	// content is a map of registered templ components for content site routes
	content map[string]templ.Component
}

func (tr *TemplRegistry) RegisterAPI(method, route string, c templ.Component) {
	if _, registered := tr.api[route]; !registered {
		tr.api[route] = make(map[string]templ.Component)
	}
	if _, registered := tr.api[route][method]; !registered {
		// Only register on the first call to this method
		tr.api[route][method] = c
	}
}

func (tr *TemplRegistry) RegisterContent(route string, c templ.Component) {
	if _, registered := tr.content[route]; !registered {
		// Only register on the first call to this method
		tr.content[route] = c
	}
}

func (tr *TemplRegistry) API(method, route string) (
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

func (tr *TemplRegistry) Content(route string) (c templ.Component, err error) {
	c, registered := tr.content[route]
	if !registered {
		return c, ErrRouteNotFound(route)
	}
	return c, nil
}

var RegisterAPI func(tr *TemplRegistry) = nil

var RegisterContent func(tr *TemplRegistry) = nil

func NewTemplRegistry() (tr *TemplRegistry) {
	tr = &TemplRegistry{
		api:     make(map[string]map[string]templ.Component),
		content: make(map[string]templ.Component),
	}

	// Register funcs are optionally defined on init
	if RegisterAPI != nil {
		RegisterAPI(tr)
	}
	if RegisterContent != nil {
		RegisterContent(tr)
	}

	return tr
}
