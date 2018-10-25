package routes

import (
	"../util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"strconv"
)

var db, err = buntdb.Open("users.db")

type user struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	Id string `json:"id"`
}

func AddUser(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	token := util.GenerateRandomString(30)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	var idString (string)
	err = db.Update(func(tx *buntdb.Tx) error {
		numUsers, err := tx.Len()
		if err != nil{
			return err;
		}
		id := numUsers + 1
		idString = strconv.Itoa(id)
		mapD := map[string]string{"name": name, "email": email, "token": token, "id": idString}
		mapB, _ := json.Marshal(mapD)
		_, _, errr := tx.Set(idString, string(mapB), nil)
		return errr
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"id": idString, "name": name, "email": email, "token": token})
	}
}

func GetUser(c *gin.Context) {
	user := c.Params.ByName("id")
	c.JSON(200, gin.H{"user": user, "message": "User get"})
}


func GetAllUsers(c *gin.Context) {
	db.View(func(tx *buntdb.Tx) error {
		numUsers, error := tx.Len()
		var users []user
		if error != nil {
			c.JSON(500, gin.H{"error": error.Error()})
		}
		for i := 0; i < numUsers; i++ {
			u, err := tx.Get(strconv.Itoa(i))
			if err != nil{
				c.JSON(500, gin.H{"error": err.Error()})
			}
			in := []byte(u)
			var raw user
			json.Unmarshal(in, &raw)
			users = append(users, raw)
		}
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"users": users})
		return err
	})
	return
}
