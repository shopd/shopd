package templrendr

import (
	"github.com/a-h/templ"
)

type Registry struct {
	// api is a map of registered templ components for api routes and methods
	api map[string]map[string]templ.Component
	// content is a map of registered templ components for content site routes
	content map[string]templ.Component
}

// Register a template component for a method and route
func (tr *Registry) Register(method, route string, c templ.Component) {
	if _, registered := tr.api[route]; !registered {
		tr.api[route] = make(map[string]templ.Component)
	}
	if _, registered := tr.api[route][method]; !registered {
		// Only register on the first call to this method
		tr.api[route][method] = c
	}
}

// RegisterContent register a content template component.
// Content components only support the GET method,
// and they all receive the same view model.
// Auth is only done on subsequent API calls,
// the content routes are the entry points for the web app
func (tr *Registry) RegisterContent(route string, c templ.Component) {
	if _, registered := tr.content[route]; !registered {
		// Only register on the first call to this method
		tr.content[route] = c
	}
}

// ByMethod returns components registered with a method and route
func (tr *Registry) ByMethod(method, route string) (
	c templ.Component, err error) {

	methods, registered := tr.api[route]
	if !registered {
		return c, ErrRouteNotFound(route)
	}

	c, registered = methods[method]
	if !registered {
		return c, ErrMethodNotSupported(method, route)
	}

	return c, nil
}

// ByRoute returns components that only support the GET method by route
func (tr *Registry) ByRoute(route string) (c templ.Component, err error) {
	c, registered := tr.content[route]
	if !registered {
		return c, ErrRouteNotFound(route)
	}
	return c, nil
}

type RegisterFunc func(tr *Registry)

func NewRegistry() (tr *Registry) {
	tr = &Registry{
		api:     make(map[string]map[string]templ.Component),
		content: make(map[string]templ.Component),
	}

	return tr
}
