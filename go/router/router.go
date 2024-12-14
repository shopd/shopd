package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shopd/shopd/go/share"
	"github.com/shopd/shopd/go/templrendr"
	"github.com/shopd/shopd/www/components"
)

// TODO Pass in services
func NewRouter() *gin.Engine {
	r := gin.Default()

	// TODO Zerolog Integration with Gin
	// https://g.co/gemini/share/70fd8e96abb5
	// r.Use(ginzerolog.Logger("gin"))

	tr := templrendr.NewRegistry()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/login", func(c *gin.Context) {
		t, err := tr.API(share.GET, "/api/login")
		if err != nil {
			// TODO Not found
			log.Error().Stack().Err(err).Msg("")
		}
		r := templrendr.New(c.Request.Context(), http.StatusOK, t)
		c.Render(http.StatusOK, r)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	r.GET("/login", func(c *gin.Context) {
		t, err := tr.Content("/login")
		if err != nil {
			// TODO Not found
			log.Error().Stack().Err(err).Msg("")
		}
		// TODO Create NewLayout constructor in templrenderer?
		l := components.Layout(t)
		r := templrendr.New(c.Request.Context(), http.StatusOK, l)
		c.Render(http.StatusOK, r)
	})

	return r
}
