package middleware

import (
	"arjunomray/eventer/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
