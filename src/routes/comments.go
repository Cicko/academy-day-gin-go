package routes

import "github.com/gin-gonic/gin"

func AddComment(c *gin.Context) {
	user := c.Params.ByName("user")
	c.JSON(200, gin.H{"message": "User created", "user": user})
}
