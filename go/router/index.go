package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/go/templrendr"
	"github.com/shopd/shopd/www/api"
	"github.com/shopd/shopd/www/content"
	"github.com/shopd/shopd/www/view"
)

func Index(c *gin.Context) {
	c.Render(http.StatusOK, templrendr.New(
		c.Request.Context(), http.StatusOK, content.Index(view.Content{})))
}

func ApiIndex(c *gin.Context) {
	c.Render(http.StatusOK, templrendr.New(
		c.Request.Context(), http.StatusOK, api.Index()))
}
