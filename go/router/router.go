package router

import (
	"github.com/gin-gonic/gin"
)

// TODO Pass in services
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	// TODO Zerolog Integration with Gin
	// https://g.co/gemini/share/70fd8e96abb5
	// r.Use(ginzerolog.Logger("gin"))

	// ...........................................................................
	// Standard routes

	// index
	r.GET("/", Index)
	r.GET("/api", ApiIndex)

	// login
	r.GET("/login", GetLogin)
	r.POST("/api/login", PostLoginAttempt)

	return r
}

// TODO Is this even useful? Rather don't match errors,
// and remove go/middleware/errorhandler.go
func ErrorMatcher(err error) (obj any, matched bool) {
	return nil, false
}
