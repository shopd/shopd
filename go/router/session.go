package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopd/shopd/go/templrendr"
	"github.com/shopd/shopd/www/api/login"
	"github.com/shopd/shopd/www/view"
)

func GetLogin(c *gin.Context) {
	rendr := templrendr.New(
		c.Request.Context(), http.StatusOK, login.Get(view.LoginGet{}))
	c.Render(http.StatusOK, rendr)
}
