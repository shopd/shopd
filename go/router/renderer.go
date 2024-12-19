package router

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

// NewRenderer constructs a gin renderer for a templ component
// https://github.com/a-h/templ/blob/main/examples/integration-gin/main.go
func NewRenderer(ctx context.Context, component templ.Component) *Renderer {
	return &Renderer{
		Ctx:       ctx,
		Component: component,
	}
}

type Renderer struct {
	Ctx       context.Context
	Component templ.Component
}

func (t Renderer) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	if t.Component != nil {
		return t.Component.Render(t.Ctx, w)
	}
	return nil
}

func (t Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
