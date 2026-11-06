package middleware

import (
	"fmt"
	"github.com/rafialfauzan/layrd/pkg/response"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery handles panics cleanly and ensures API always returns a consistent JSON error
func Recovery(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Printf("[Panic Recovered] Error: %v\nStack trace:\n%s\n", err, stack)

				var errDetail any = "An unexpected server error occurred."
				if isDebug {
					errDetail = fmt.Sprintf("%v", err)
				}

				response.InternalError(c, "Internal Server Error", errDetail)
				c.Abort()
			}
		}()
		c.Next()
	}
}
