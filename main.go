package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	PicSRC string `json:"pic"`
	Status bool   `json:"status"`
}

var users []User = []User{
	{ID: "1", Name: "John Doe", PicSRC: "mockup-1", Status: false},
	{ID: "2", Name: "Jane Doe", PicSRC: "mockup-2", Status: false},
	{ID: "3", Name: "Sponge Bob", PicSRC: "mockup-3", Status: false},
}

func getUsers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, users)
}

func registerUser(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserByID(id string) (*User, error) {
	for i, user := range users {
		if user.ID == id {
			return &users[i], nil
		}
	}
	return nil, errors.New("User not found")
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := getUserByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func statusToggle(c *gin.Context) {
	id := c.Param("id")
	user, err := getUserByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	user.Status = !user.Status
	c.IndentedJSON(http.StatusOK, user)
}

func removeUserByID(id string) error {
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}
	return errors.New("User not found")
}

func removeUser(c *gin.Context) {
	id := c.Param("id")
	err := removeUserByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func main() {
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/register", registerUser)
	router.PATCH("/user/:id", statusToggle)
	router.DELETE("/user/:id", removeUser)
	router.GET("/user/:id", getUser)
	router.Run("localhost:4000")
}
