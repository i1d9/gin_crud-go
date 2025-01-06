package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/i1d9/gin_crud-go/models"
	"github.com/jackc/pgx/v5/pgxpool"
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "something went wrong"})
			c.Abort()
			return
		}

		now := time.Now().UTC()

		session_expiry := session.Expires_At.UTC()

		if now.Before(session_expiry) {
			c.Set("session", session)

			c.Next()

		} else {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "message": "invalid token"})
			c.Abort()
			return
		}

	}
}
