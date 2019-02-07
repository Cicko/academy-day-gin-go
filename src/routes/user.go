package routes

import (
	"../util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/buntdb"
	"strconv"
)

var UserDb, _ = buntdb.Open("users.db")

type user struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	Id string `json:"id"`
}

func AddUser(c *gin.Context) {
	// We can not use PostForm because in the Tests the user is created providing RAW data
	// name := c.PostForm("name")
	// email := c.PostForm("email")
	rawData, _ := c.GetRawData()
	reqForm := &user{}
	err := json.Unmarshal([]byte(rawData), reqForm)
	token := util.GenerateRandomString(30)

	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	var idString (string)
	var id (int)
	err = UserDb.Update(func(tx *buntdb.Tx) error {
		err := tx.CreateIndex("token", "token")
		if err != nil{
			fmt.Errorf(err.Error())
		}
		numUsers, err := tx.Len()
		if err != nil{
			return err;
		}
		id = numUsers + 1
		idString = strconv.Itoa(id)
		mapD := map[string]string{"name": reqForm.Name, "email": reqForm.Email, "token": token, "id": idString}
		mapB, _ := json.Marshal(mapD)
		_, _, errr := tx.Set(idString, string(mapB), nil)
		return errr
	})
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"id": id, "name": reqForm.Name, "email": reqForm.Email, "token": token})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")

	UserDb.View(func(tx *buntdb.Tx) error {
		user, err := tx.Get(id)
		if err != nil{
			c.JSON(404, gin.H{"error": "User doesn't exist"})
			return err
		}
		c.JSON(200, ReformatUser(user))
		return err
	})
}


func GetAllUsers(c *gin.Context) {
	UserDb.View(func(tx *buntdb.Tx) error {
		numUsers, error := tx.Len()
		users := make([]user, 0)
		if error != nil {
			c.JSON(500, gin.H{"error": error.Error()})
			return error
		}
		for i := 1; i <= numUsers; i++ {
			u, err := tx.Get(strconv.Itoa(i))
			if err == nil {
				users = append(users, ReformatUser(u))
			}
		}
		c.JSON(200, gin.H{"users": users})
		return err
	})
	return
}

func EditUser(c *gin.Context) {
	id := c.Params.ByName("id")
	rawData, _ := c.GetRawData()
	reqForm := &user{}
	err := json.Unmarshal([]byte(rawData), reqForm)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	name := reqForm.Name
	email := reqForm.Email
	UserDb.Update(func(tx *buntdb.Tx) error {
		u, err := tx.Get(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "User doesn't exist"})
			return err
		}
		reformattedUser := ReformatUser(u)
		mapD := map[string]interface{}{
			"token": reformattedUser.Token,
			"id": reformattedUser.Id,
		}
		if name != "" {
			mapD["name"] = name
		}
		if email != "" {
			mapD["email"] = email
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

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	UserDb.Update(func(tx *buntdb.Tx) error {
		// tx.Set(id, "cucu", &buntdb.SetOptions{Expires:true, TTL:time.Second})
		_, err := tx.Delete(id)
		if err != nil{
			c.JSON(500, gin.H{"error": err.Error()})
			return err
		}
		c.JSON(200, nil)
		return err
	})
}

func DeleteAllUsers(c *gin.Context) {
	UserDb.Update(func(tx *buntdb.Tx) error {
		error := tx.DeleteAll()
		if error != nil{
			c.JSON(500, gin.H{"error": error.Error()})
			return error
		}
		c.JSON(200, nil)
		return nil
	})
}

// Utils
func ReformatUser(u string) user {
	in := []byte(u)
	var raw user
	json.Unmarshal(in, &raw)
	return raw
}


func GetUserByToken(token string) (user, error) {
	var wantedUser user
	err := UserDb.View(func(tx *buntdb.Tx) error {
		numUsers, error := tx.Len()
		if error != nil {
			fmt.Errorf(err.Error())
			return error
		}
		for i := 1; i <= numUsers; i++ {
			u, err := tx.Get(strconv.Itoa(i))
			if err != nil{
				fmt.Errorf(err.Error())
				return error
			}
			reformattedUser := ReformatUser(u)
			fmt.Println(reformattedUser)
			if token == reformattedUser.Token {
				wantedUser = reformattedUser
			}
		}
		if err != nil{
			fmt.Errorf(err.Error())
			return err
		}
		return nil
	})
	return wantedUser, err
}
