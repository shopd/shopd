package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/www/api"
	"github.com/shopd/shopd/www/content"
)

func (h *RouteHandler) Index(c *gin.Context) {
	c.Render(http.StatusOK, h.Content(c.Request, content.Index))
}

func (h *RouteHandler) ApiIndex(c *gin.Context) {
	c.Render(http.StatusOK, h.Content(c.Request, api.Index))
}
