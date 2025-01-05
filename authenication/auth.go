package authenication

import (
	"github.com/gin-gonic/gin"
	"github.com/i1d9/gin_crud-go/middleware"
	"github.com/i1d9/gin_crud-go/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

func Login(c *gin.Context, pool *pgxpool.Pool) {

	// Get user details
	identifier := c.PostForm("identifier")

	// Get password
	password := c.PostForm("password")

	user, err := models.GetAuthUser(pool, identifier)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid identifier/password are invalid",
		})
		return
	}

	ok, err := argon2.VerifyEncoded([]byte(password), []byte(user.Password))
	if err != nil {

		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid identifier/password are invalid",
		})
		return
	}

	if ok {

		session_id, session_creation_err := models.CreateSession(pool, user.ID)
		if session_creation_err != nil {
			c.JSON(400, gin.H{

				"error": session_creation_err.Error(),
			})
			return
		}

		session, session_fetch_err := models.GetSessionByID(pool, session_id)

		if session_creation_err != nil || session_fetch_err != nil {
			c.JSON(400, gin.H{

				"error": session_fetch_err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"token_type":   "Bearer",
			"expires_in":   3600,
			"access_token": session.Token,
		})
		return

	} else {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid identifier/password are invalid",
		})
		return
	}

}

func Register(c *gin.Context, pool *pgxpool.Pool) {

	// Get user details
	username := c.PostForm("username")
	gender := c.PostForm("gender")
	email := c.PostForm("email")
	mobile_number := c.PostForm("mobile_number")
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	surname := c.PostForm("surname")
	// Get password
	password := c.PostForm("password")

	err := models.CreateUser(pool, first_name, last_name, surname, email, username, mobile_number, gender, password)

	if err != nil {
		c.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid parameter were provided",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Account created successfully",
	})

}

func Logout(c *gin.Context, pool *pgxpool.Pool) {

	access_token := middleware.ExtractAccessToken(c)

	session, session_fetch_err := models.GetSessionByTokenAndTokenType(pool, access_token, "access_token")

	if session_fetch_err != nil {
		c.JSON(400, gin.H{
			"error":   session_fetch_err.Error(),
			"message": "session not found",
		})
		return
	}

	rows, session_delete_err := models.DeleteSession(pool, session.ID)

	if session_delete_err != nil || rows == 0 {
		c.JSON(400, gin.H{
			"error":   session_delete_err.Error(),
			"message": "Invalid identifier/password are invalid",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully logged out",
	})

}
