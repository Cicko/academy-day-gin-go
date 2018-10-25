package routes

import (
	"../util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"time"
)

var db, err = buntdb.Open("users.db")

type User struct {
	Name string
	Email string
	Token string
}

func AddUser(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	token := util.GenerateRandomString(30)

	user := &User{Name: name, Email: email, Token: token}

	b, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	t := time.Now().String()
	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(t, string(b), nil)
		return err
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"id": t, "name": name, "email": email, "token": token})
	}
}

func GetUser(c *gin.Context) {
	user := c.Params.ByName("id")
	c.JSON(200, gin.H{"user": user, "message": "User get"})
}
