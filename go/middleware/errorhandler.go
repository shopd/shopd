package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMatcher interface {
	// Match the error and return the object to be rendered,
	// otherwise return false if the error is not matched
	Match(err error) (obj any, matched bool)
}

type ErrorHandlerParams struct {
	// TODO Toggle output JSON or HTML
	Matchers []ErrorMatcher
}

// NewErrorHandler returns error handler middleware for use with gin
// https://stackoverflow.com/a/69948929/639133
func ErrorHandler(params ErrorHandlerParams) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			if len(c.Errors) > 1 {
				// Convention is to not allow multiple errors
				c.JSON(http.StatusInternalServerError, "Multiple errors")
				return
			}

			// Only match on the first error
			err := c.Errors[1]
			if c.Writer.Status() == http.StatusOK {
				// Default status code but an error was set
				c.Status(http.StatusInternalServerError)
			}

			// Try to match the error
			for _, matcher := range params.Matchers {
				obj, matched := matcher.Match(err.Err)
				if matched {
					// Error is rendered on the first match,
					// assume the status code has already been set and don't override.
					// The handler may perform additional logic depending on the error
					c.JSON(-1, obj)
					return
				}
			}

			// Unmatched error
			c.JSON(-1, err.Error())
		}
	}
}
