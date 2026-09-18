package middleware

import (
	"strings"

	"arjunomray/eventer/api"
	"arjunomray/eventer/models"
	"arjunomray/eventer/utils"

	"github.com/gin-gonic/gin"
)

const claimsKey = "claims"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || token == "" {
			api.Unauthorized(c, "Invalid or missing authorization header")
			return
		}

		claims, err := utils.VerifyJWT(token)
		if err != nil {
			api.Unauthorized(c, "Invalid or expired token")
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

// RequireAdmin allows only administrators through. It must run after
// AuthMiddleware.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := ClaimsFrom(c)
		if !ok || claims.Role != models.RoleAdmin {
			api.Forbidden(c, "Administrator access required")
			return
		}
		c.Next()
	}
}

// ClaimsFrom returns the authenticated user's claims, if present.
func ClaimsFrom(c *gin.Context) (*utils.Claims, bool) {
	value, exists := c.Get(claimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := value.(*utils.Claims)
	return claims, ok
}
