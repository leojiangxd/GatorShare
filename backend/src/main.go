package main

import (
	"fmt"

	"backend/src/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/posts", api.GetPosts)

	fmt.Println("Starting up server...")
	router.Run("localhost:8080")
}