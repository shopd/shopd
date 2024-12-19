package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO Pass in services
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	// TODO Zerolog Integration with Gin
	// https://g.co/gemini/share/70fd8e96abb5
	// r.Use(ginzerolog.Logger("gin"))

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// r.GET("/api/login", func(c *gin.Context) {
	// 	t, err := tr.API(share.GET, "/api/login")
	// 	if err != nil {
	// 		// TODO Not found
	// 		log.Error().Stack().Err(err).Msg("")
	// 	}
	// 	r := templrendr.New(c.Request.Context(), http.StatusOK, t)
	// 	c.Render(http.StatusOK, r)
	// })

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	return r
}

// TODO Is this even useful? Rather don't match errors,
// and remove go/middleware/errorhandler.go
func ErrorMatcher(err error) (obj any, matched bool) {
	return nil, false
}
