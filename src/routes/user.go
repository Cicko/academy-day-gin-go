package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
)

var db, err = buntdb.Open("users.db")


type User struct {
	name string
	email string
}

func AddUser(c *gin.Context) {
	name := c.Params.ByName("name")
	email := c.Params.ByName("email")
	fmt.Println(name)
	fmt.Println(email)
	user := &User{name: "Ruda", email: "rudo"}
	b, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user)
	c.JSON(200, gin.H{"message": "User created", "user": string(b)})
	/*
	err := db.Update(func(tx *buntdb.Tx) error {
		t := time.Now().String()
		_, _, err := tx.Set(t, "myvalue", nil)
		return err
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error(),})
	} else {
		c.JSON(200, gin.H{"message": "User created", "user": user})
	}
	*/
}

func GetUser(c *gin.Context) {
	user := c.Params.ByName("id")
	c.JSON(200, gin.H{"user": user, "message": "User get"})
}
