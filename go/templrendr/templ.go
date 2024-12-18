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

// Register is used to register templates for the Hypermedia API
var Register func(tr *Registry) = nil

// RegisterContent is used to register content templates
var RegisterContent func(tr *Registry) = nil

func NewRegistry() (tr *Registry) {
	tr = &Registry{
		api:     make(map[string]map[string]templ.Component),
		content: make(map[string]templ.Component),
	}

	// Defined in the generated init func in the router package...
	if Register != nil {
		Register(tr)
	}
	if RegisterContent != nil {
		RegisterContent(tr)
	}

	return tr
}
