package middleware

import (
	"strings"

	"github.com/rafialfauzan/layrd/pkg/response"
	"github.com/rafialfauzan/layrd/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
	ContextRole     = "role"
	ContextClaims   = "jwt_claims"
)

// JWTAuth verifies incoming bearer token and sets authenticated user claims into gin.Context
func JWTAuth(tokenService *utils.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Invalid Authorization header format. Expected 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := tokenService.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Inject identity into Gin context
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextRole, claims.Role)
		c.Set(ContextClaims, claims)

		c.Next()
	}
}

// RequireRole checks whether the authenticated user has one of the required roles
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[strings.ToLower(r)] = true
	}

	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextRole)
		if !exists {
			response.Forbidden(c, "User role not identified")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok || !allowed[strings.ToLower(userRole)] {
			response.Forbidden(c, "You do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
