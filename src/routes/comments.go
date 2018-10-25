package routes

import (
	"github.com/gin-gonic/gin"
)

type comment struct {
	Message string `json:"message"`
	PostId string `json:"postId"`
	Author string `json:"author"`
	Id string `json:"id"`
}

func AddComment(c *gin.Context) {
	user := c.Params.ByName("user")
	c.JSON(200, gin.H{"message": "User created", "user": user})
}
