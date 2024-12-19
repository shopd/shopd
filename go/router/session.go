package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/www/api/login"
	content "github.com/shopd/shopd/www/content/login"
	"github.com/shopd/shopd/www/view"
)

func (h *RouteHandler) GetLogin(c *gin.Context) {
	c.Render(http.StatusOK, h.Content(c.Request, content.Index(view.Content{})))
}

func (h *RouteHandler) PostLoginAttempt(c *gin.Context) {
	c.Render(http.StatusOK, h.Template(c.Request, login.Post(view.LoginPost{})))
}
