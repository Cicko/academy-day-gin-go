package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"strconv"
)

var PostsDb, err = buntdb.Open("posts.db")

type post struct {
	Message string `json:"message"`
	Comments []comment `json:"comments"`
	Author user `json:"author"`
	Id string `json:"id"`
}

func AddPost(c *gin.Context) {
	message := c.PostForm("message")
	token := c.GetHeader("token")
	var author user

	err := UserDb.View(func(tx *buntdb.Tx) error {
		author = GetUserByToken(token)
		return nil
	})


	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	var idString (string)
	err = PostsDb.Update(func(tx *buntdb.Tx) error {
		numPosts, err := tx.Len()
		if err != nil{
			return err;
		}
		id := numPosts + 1
		idString = strconv.Itoa(id)
		post := &post{Message: message, Author: author, Id: idString}
		mapB, _ := json.Marshal(post)
		_, _, errr := tx.Set(idString, string(mapB), nil)
		return errr
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"id": idString, "message": message, "author": author})
	}
}

