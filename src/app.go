package main

import (
	"./routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	Users := router.Group("/api/users")
	{
		Users.POST("/", routes.AddUser)
		Users.GET("/:id", routes.GetUser)
		Users.GET("/", routes.GetAllUsers)
		Users.PUT("/:id", routes.EditUser)
		Users.DELETE("/:id", routes.DeleteUser)
	}
	Posts := router.Group("api/posts")
	{
		Posts.POST("/", routes.AddPost)
		Posts.GET("/:id", routes.GetPost)
		Posts.PUT("/:id", routes.EditPost)
		Posts.DELETE("/:id", routes.DeletePost)
		Posts.GET("/", routes.ShowPosts)
	}
	Comments := router.Group("api/posts")
	{
		Comments.POST("/:postId/comments", routes.AddComment)
		// Comments.GET("/:postId/comments/:id", routes.GetComment)
	}
	router.Run()
}
