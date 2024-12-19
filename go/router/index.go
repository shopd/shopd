package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/www/api"
	"github.com/shopd/shopd/www/content"
	"github.com/shopd/shopd/www/view"
)

func (h *RouteHandler) Index(c *gin.Context) {
	c.Render(http.StatusOK, h.Content(c.Request, content.Index(view.Content{})))
}

func (h *RouteHandler) ApiIndex(c *gin.Context) {
	c.Render(http.StatusOK, h.Content(c.Request, api.Index()))
}
