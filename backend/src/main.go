package main

import (
	"fmt"
	"log"

	"backend/src/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/posts", api.GetPosts)

	fmt.Println("Starting up server...")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}