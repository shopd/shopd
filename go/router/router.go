package router

import (
	"net/http"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd-proto/go/share"
	"github.com/shopd/shopd/go/config"
	"github.com/shopd/shopd/www/components"
	"github.com/shopd/shopd/www/view"
)

// ContentFunc is the signature for all content templates
type ContentFunc func(model view.Content) templ.Component

type RouteHandler struct {
	model view.Content
}

// Template renders a templ component
func (h *RouteHandler) Template(
	r *http.Request, template templ.Component) *Renderer {
	return NewRenderer(r.Context(), template)
}

// Content renders a templ component with layout
func (h *RouteHandler) Content(
	r *http.Request, content ContentFunc) *Renderer {
	return NewRenderer(r.Context(), components.Layout(h.model, content(h.model)))
}

// TODO Pass in services
func NewRouter(conf *config.Config) *gin.Engine {
	h := RouteHandler{}
	h.model = view.NewContent(view.ContentParams{
		BaseURL:      "https://localhost:8443/",  // TODO Use config
		DomainConfig: share.DomainConfigExport{}, // TODO Domain config
	})

	r := gin.Default()
	r.Use(gin.Recovery())

	// TODO Zerolog Integration with Gin
	// https://g.co/gemini/share/70fd8e96abb5
	// r.Use(ginzerolog.Logger("gin"))

	// ...........................................................................
	// Standard routes

	// index
	r.GET("/", h.Index)
	r.GET("/api", h.ApiIndex)

	// login
	r.GET("/login", h.GetLogin)
	r.POST("/api/login", h.PostLoginAttempt)

	// static
	staticRoot := filepath.Join(conf.Dir(), "www", "static")
	r.Static("/s", staticRoot)

	return r
}

// TODO Is this even useful? Rather don't match errors,
// and remove go/middleware/errorhandler.go
func ErrorMatcher(err error) (obj any, matched bool) {
	return nil, false
}
