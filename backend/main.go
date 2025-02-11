package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"testsite.com/sr1/models"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	//API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("user", getUsers)
		v1.GET("user/:id", getUserById)
		v1.POST("user", addUser)
		v1.PUT("user/:id", updateUser)
		v1.DELETE("user/:id", deleteUser)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()

}

func getUsers(c *gin.Context) {

	users, err := models.GetUsers(10)
	checkErr(err)

	if users == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func getUserById(c *gin.Context) {

	id := c.Param("id")

	user, err := models.GetUserById(id)
	checkErr(err)

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func addUser(c *gin.Context) {

	var json models.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddUser(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updateUser(c *gin.Context) {

	var json models.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateUser(json, userId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteUser(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteUser(userId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
