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

func ConnectDatabase() error {
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
	err := ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	//API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", index)
		v1.GET("member", getMembers)
		v1.GET("member/:id", getMemberById)
		v1.POST("member", addMember)
		v1.PUT("member/:id", updateMember)
		v1.DELETE("member/:id", deleteMember)
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

	result := db.Where("id = ?", id).First(&member)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": member})
	}
}

func addMember(c *gin.Context) {

	var newMember models.Member

	if err := c.ShouldBindJSON(&newMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&newMember)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

func updateMember(c *gin.Context) {

	id := c.Param("id")

	var member models.Member

	if err := db.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&member)
	c.JSON(http.StatusOK, gin.H{"data": member})
}

func deleteMember(c *gin.Context) {

	id := c.Param("id")

	var member models.Member

	if err := db.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	db.Delete(&member)
	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

func options(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "options Called"})

}
