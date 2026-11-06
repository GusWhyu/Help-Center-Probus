package middleware

import (
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/response"
)

// Recovery handles panics cleanly and ensures API always returns a consistent JSON error
func Recovery(_ bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Printf("[Panic Recovered] Error: %v\nStack trace:\n%s\n", err, stack)

				response.InternalError(c, "Internal Server Error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
