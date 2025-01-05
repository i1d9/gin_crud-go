package authenication

import "github.com/gin-gonic/gin"

func Login(c *gin.Context) {

	// Get user details
	identifier := c.PostForm("identifier")

	// Get password
	password := c.PostForm("password")

	c.JSON(200, gin.H{
		"password":   password,
		"identifier": identifier,
	})
}

func Register(c *gin.Context) {
	// Get user details
	username := c.PostForm("username")
	// email := c.PostForm("email")
	// mobile_number := c.PostForm("mobile_number")
	// first_name := c.PostForm("first_name")
	// last_name := c.PostForm("last_name")
	// surname := c.PostForm("surname")
	// Get password
	password := c.PostForm("password")

	c.JSON(200, gin.H{
		"password":   password,
		"identifier": username,
	})
}

func Logout(c *gin.Context) {
	// Get user details
	username := c.PostForm("username")
	// email := c.PostForm("email")
	// mobile_number := c.PostForm("mobile_number")
	// first_name := c.PostForm("first_name")
	// last_name := c.PostForm("last_name")
	// surname := c.PostForm("surname")
	// Get password
	password := c.PostForm("password")

	c.JSON(200, gin.H{
		"password":   password,
		"identifier": username,
	})
}
