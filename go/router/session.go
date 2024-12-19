package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/go/templrendr"
	"github.com/shopd/shopd/www/api/login"
	"github.com/shopd/shopd/www/components"
	content "github.com/shopd/shopd/www/content/login"
	"github.com/shopd/shopd/www/view"
)

func GetLogin(c *gin.Context) {
	c.Render(http.StatusOK, templrendr.New(
		c.Request.Context(), http.StatusOK, components.Layout(
			content.Index(view.Content{}))))
}

func PostLoginAttempt(c *gin.Context) {
	c.Render(http.StatusOK, templrendr.New(
		c.Request.Context(), http.StatusOK, login.Post(view.LoginPost{})))
}
