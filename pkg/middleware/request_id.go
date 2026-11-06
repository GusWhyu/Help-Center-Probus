package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const HeaderXRequestID = "X-Request-ID"

// RequestID attaches or propagates a unique request identifier
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			bytes := make([]byte, 12)
			_, _ = rand.Read(bytes)
			reqID = hex.EncodeToString(bytes)
		}

		c.Set(HeaderXRequestID, reqID)
		c.Header(HeaderXRequestID, reqID)
		c.Next()
	}
}
