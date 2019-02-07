package main

import (
	"fmt"
	"log"
	"time"
	"./routes"
	"github.com/gin-gonic/gin"
	"github.com/radovskyb/watcher"
)

func configureGin() {
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

func main() {
	configureGin()
	w := watcher.New()

	go func() {
			for {
				select {
				case event := <-w.Event:
					fmt.Println(event) // Print the event's info.
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
	}()

	if err := w.Add("."); err != nil {
		log.Fatalln(err)
	}

	if err := w.AddRecursive("./routes"); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	// Trigger 2 events after watcher started.
	go func() {
		w.Wait()
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
