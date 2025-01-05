package gincrud

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/i1d9/gin_crud-go/authenication"
)



func main() {

	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/login", authenication.Login)
		auth.POST("/register", authenication.Register)
		auth.POST("/logout", authenication.Logout)

	}


	user := router.Group("/user")
	{
		user.POST("/search", authenication.Login)
		user.GET("/profile", authenication.Register)
		
	}

	
	router.GET("/search", func(c *gin.Context) {
		queryValue := c.Query("term") // Extract query parameter "term"
		c.JSON(200, gin.H{
			"term": queryValue,
		})
	})

	
	

	router.Run(":8080")
}
