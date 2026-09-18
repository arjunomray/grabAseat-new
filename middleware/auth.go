package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"arjunomray/eventer/api"
	"arjunomray/eventer/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.Unauthorized(c, "Authorization header is required")
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			api.Unauthorized(c, "Token is required")
			return
		}
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			api.Unauthorized(c, "Invalid or expired token")
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
