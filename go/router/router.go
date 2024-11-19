package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO Pass in services
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// TODO Zerolog Integration with Gin
	// https://g.co/gemini/share/70fd8e96abb5

	// TODO Templ integration
	// https://templ.guide/integrations/web-frameworks/
	// https://github.com/a-h/templ/blob/main/examples/integration-gin/main.go

	return r
}
