package middleware

import (
	"arjunomray/eventer/api"
	"arjunomray/eventer/utils"
	"strings"

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

// ClaimsFrom returns the authenticated user's claims, if present.
func ClaimsFrom(c *gin.Context) (*utils.Claims, bool) {
	value, exists := c.Get(claimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := value.(*utils.Claims)
	return claims, ok
}
