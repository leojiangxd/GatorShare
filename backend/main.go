package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"gshare.com/platform/models"
)

var db *gorm.DB

func connectDatabase() error {
	var err error
	db, err = gorm.Open(sqlite.Open("local.db"), &gorm.Config{})
	if err != nil {
		panic("Error connecting/creating the sqlite db")
	}
	db.AutoMigrate(&models.Member{})
	return err
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := connectDatabase()
	checkErr(err)

	r := gin.Default()

	//API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", index)
		v1.GET("member", getMembers)
		//v1.GET("member/:id", getMemberById)
		v1.GET("member/:username", getMemberByUsername)
		v1.POST("register", register)
		v1.PUT("member/:username", updateMember)
		v1.DELETE("member/:username", deleteMember)
		v1.POST("login", login)
		v1.POST("logout/:username", logout)
		v1.OPTIONS("member", options)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}

func index(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "index Called"})

}

func getMembers(c *gin.Context) {

	var members []models.Member

	result := db.Limit(10).Find(&members)

	checkErr(result.Error)

	if members == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": members})
	}
}

func getMemberById(c *gin.Context) {

	id := c.Param("id")

	var member models.Member

	result := db.First(&member, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": member})
	}
}

func getMemberByUsername(c *gin.Context) {

	username := c.Param("username")

	var member models.Member

	result := db.First(&member, "username = ?", username)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": member})
	}
}

func register(c *gin.Context) {

	var newMember models.Member

	//Check that the request is in the correct format
	if err := c.ShouldBindJSON(&newMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if there is already a user with that username (could implement frontend requirements for this without using http requests)
	if err := db.First(&newMember, "username = ?", newMember.Username).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	//Hash the password using bcrypt
	newMember.Password, _ = hashPassword(newMember.Password)

	//Add to database
	result := db.Create(&newMember)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

func login(c *gin.Context) {

	var loginInfo models.Member

	//Bind the username and password into the member struct
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := loginInfo.Username
	password := loginInfo.Password

	//Check username and password match
	var member models.Member
	result := db.First(&member, "username = ?", username)
	if result == nil || !checkPasswordHash(password, member.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password"})
		return
	}

	//Generate tokens and pass to user cookies
	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.SetCookie("csrf_token", csrfToken, 3600, "/", "localhost", false, false)

	//Store tokens in the database
	db.Model(&member).Where("username = ?", username).Updates(map[string]any{"session_token": sessionToken, "csrf_token": csrfToken})
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func logout(c *gin.Context) {

	//Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Clear cookies
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("csrf_token", "", -1, "/", "localhost", false, false)

	username := c.Param("username")

	var member models.Member

	//Check that user exists
	if err := db.First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such user" + username})
		return
	}

	//Clear tokens from db
	member.SessionToken = ""
	member.CSRFToken = ""
	db.Model(&member).Where("username = ?", username).Updates(map[string]any{"session_token": "", "csrf_token": ""})
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func updateMember(c *gin.Context) {

	//Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := c.Param("username")

	var member models.Member

	//Check that user exists
	if err := db.First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No such user" + username})
		return
	}

	//Bind updated user
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if you have to hash new password
	if member.Password != "" {
		var err error
		if member.Password, err = hashPassword(member.Password); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
	}

	db.Model(&member).Where("username = ?", username).Updates(models.Member{Email: member.Email, Username: member.Username, Password: member.Password})
	c.JSON(http.StatusOK, gin.H{"data": member})
}

func deleteMember(c *gin.Context) {

	//Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	username := c.Param("username")

	var member models.Member

	if err := db.First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&member)
	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

func options(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "options Called"})

}
