package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func VerifyAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" && !strings.HasPrefix(authHeader, "Bearer "){
		c.JSON(401, gin.H{"error": "Authorization header is required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")


	

	c.Next()

}
