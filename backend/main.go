package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"testsite.com/sr1/models"
)

var db gorm.DB

func ConnectDatabase() error {

	var err error
	db, err := gorm.Open(sqlite.Open("names.db"), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&models.Person{})
	return nil
}

func main() {
	err := ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	//API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("person", getPersons)
		v1.GET("person/:id", getPersonById)
		v1.POST("person", addPerson)
		v1.PUT("person/:id", updatePerson)
		v1.DELETE("person/:id", deletePerson)
		v1.OPTIONS("person", options)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()

}

func getPersons(c *gin.Context) {

	var persons []models.Person

	result := db.Limit(10).Find(&persons)
	checkErr(result.Error)

	if persons == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": persons})
	}
}

func getPersonById(c *gin.Context) {

	id := c.Param("id")

	var person models.Person

	result := db.Where("id = ", id).First(&person)
	checkErr(result.Error)

	// if the name is blank we can assume nothing is found
	if person.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": person})
	}
}

func addPerson(c *gin.Context) {

	var newPerson models.Person

	if err := c.ShouldBindJSON(&newPerson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Create(&newPerson)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

func updatePerson(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "updatePerson Called"})
}

func deletePerson(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "deletePerson " + id + " Called"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "options Called"})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
