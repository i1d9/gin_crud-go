package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/i1d9/gin_crud-go/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"strings"
)

func ExtractAccessToken(c *gin.Context) string {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" && !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token
}

func VerifyAccessToken(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := ExtractAccessToken(c)

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		session, err := models.FindSessionbyToken(pool, token)

		if err != nil {
			log.Fatalf("Error searching for users: %v\n", err)
		}

		c.Set("session", session)

		c.Next()
	}
}
