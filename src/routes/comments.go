package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"strconv"
	"time"
)

var CommentsDb, _ = buntdb.Open("comments.db")

type comment struct {
	Message string `json:"message"`
	PostId string `json:"postId"`
	Id string `json:"id"`
}

func AddComment(c *gin.Context) {
	message := c.PostForm("message")
	token := c.GetHeader("token")
	postId := c.Params.ByName("postId")

	// Check if user exists
	err := UserDb.View(func(tx *buntdb.Tx) error {
		_, err = GetUserByToken(token)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if post exists and update it with the new comment
	err = PostsDb.Update(func(tx *buntdb.Tx) error {
		post, err := tx.Get(postId)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		reformattedPost := ReformatPost(post)
		reformattedPost.Comments = append(reformattedPost.Comments, message)

		mapD := map[string]interface{}{
			"comments": reformattedPost.Comments,
			"id": reformattedPost.Id,
			"author": reformattedPost.Author,
		}
		mapB, _ := json.Marshal(mapD)
		_, _, errr := tx.Set(postId, string(mapB), nil)
		if errr != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var idString (string)
	err = CommentsDb.Update(func(tx *buntdb.Tx) error {
		numComments, err := tx.Len()
		if err != nil{
			return err;
		}
		id := numComments + 1
		idString = strconv.Itoa(id)
		post := &comment{Message: message, Id: idString, PostId: postId}
		mapB, _ := json.Marshal(post)
		_, _, errr := tx.Set(idString, string(mapB), nil)
		return errr
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"id": idString, "message": message})
	}
}


func GetComment(c *gin.Context) {
	id := c.Params.ByName("id")
	postId := c.Params.ByName("postId")
	token := c.GetHeader("token")

	// Check if user exists
	err := UserDb.View(func(tx *buntdb.Tx) error {
		_, err = GetUserByToken(token)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = CheckPostAuthor(token, postId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	CommentsDb.View(func(tx *buntdb.Tx) error {
		comment, err := tx.Get(id)
		if err != nil{
			c.JSON(404, gin.H{"error": "Comment doesn't exist"})
			return err
		}
		rc := ReformatComment(comment)
		c.JSON(200, gin.H{"postId": rc.PostId, "message": rc.Message})
		return err
	})
}

func EditComment(c *gin.Context) {
	id := c.Params.ByName("id")
	message := c.PostForm("message")
	token := c.GetHeader("token")
	err := CheckPostAuthor(token, id)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
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

func DeleteComment(c *gin.Context) {
	id := c.Params.ByName("id")
	token := c.GetHeader("token")
	err := CheckPostAuthor(token, id)
	if err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}
	PostsDb.Update(func(tx *buntdb.Tx) error {
		tx.Set(id, "cucu", &buntdb.SetOptions{Expires:true, TTL:time.Second})
		return nil
	})
}

func GetComments(c *gin.Context) {
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
func ReformatComment(p string) comment {
	in := []byte(p)
	var raw comment
	json.Unmarshal(in, &raw)
	return raw
}