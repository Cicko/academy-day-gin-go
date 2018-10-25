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
		// Users.PUT("/:id", editUser)
		// Users.DELETE("/:id", deleteUser)
	}
	router.Run()
}
