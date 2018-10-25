package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"strconv"
	"time"
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


func GetPost(c *gin.Context) {
	id := c.Params.ByName("id")

	PostsDb.View(func(tx *buntdb.Tx) error {
		post, err := tx.Get(id)
		if err != nil{
			c.JSON(404, gin.H{"error": "Post doesn't exist"})
			return err
		}
		c.JSON(200, gin.H{"post": ReformatPost(post)})
		return err
	})
}

func EditPost(c *gin.Context) {
	id := c.Params.ByName("id")
	message := c.PostForm("message")
	PostsDb.Update(func(tx *buntdb.Tx) error {
		p, err := tx.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "Post doesn't exist"})
			return err
		}
		reformattedPost := ReformatPost(p)
		mapD := map[string]interface{}{
			"comments": reformattedPost.Comments,
			"id": reformattedPost.Id,
			"author": reformattedPost.Author,
		}
		if message != "" {
			mapD["message"] = message
		}
		mapB, _ := json.Marshal(mapD)
		_, _, errr := tx.Set(id, string(mapB), nil)
		if errr != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		c.JSON(200, gin.H(mapD))
		return nil
	})
}

func DeletePost(c *gin.Context) {
	id := c.Params.ByName("id")
	PostsDb.Update(func(tx *buntdb.Tx) error {
		tx.Set(id, "cucu", &buntdb.SetOptions{Expires:true, TTL:time.Second})
		return nil
	})
}

func ShowPosts(c *gin.Context) {
	PostsDb.View(func(tx *buntdb.Tx) error {
		numUsers, error := tx.Len()
		var posts []post
		if error != nil {
			c.JSON(500, gin.H{"error": error.Error()})
			return error
		}
		for i := 1; i <= numUsers; i++ {
			p, err := tx.Get(strconv.Itoa(i))
			if err != nil{
				c.JSON(500, gin.H{"error": err.Error()})
			}
			posts = append(posts, ReformatPost(p))
		}
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		c.JSON(200, gin.H{"posts": posts})
		return err
	})
	return
}

// Utils
func ReformatPost(p string) post {
	in := []byte(p)
	var raw post
	json.Unmarshal(in, &raw)
	return raw
}