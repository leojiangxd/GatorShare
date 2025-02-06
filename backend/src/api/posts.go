package api

import (
	// "database/sql"
	"fmt"
	"net/http"
	// "os"

	"github.com/gin-gonic/gin"
	// "github.com/mattn/go-sqlite3"
	// "github.com/joho/godotenv"
)

type Post struct {
	PostID		int		`json:"post_id"`
	UserID		int		`json:"user_id"`
	Title		string	`json:"title"`
	Date		string	`json:"date"`
	Content		string	`json:"content"`
	Likes		int		`json:"likes"`
	Dislikes	int		`json:"dislikes"`
}

// example posts
var example_posts = []Post{
	{PostID: 1, UserID: 1, Title: "Post Title", Date: "2025-01-01 00:00:00", Content: "Lorem Ipsum", Likes: 101, Dislikes: 100},
	{PostID: 2, UserID: 2, Title: "Example Post", Date: "2025-01-01 00:00:00", Content: "Lorem Ipsum", Likes: 50, Dislikes: 0},
}

func GetPosts(c *gin.Context) {
	fmt.Println("Getting posts...")
	c.IndentedJSON(http.StatusOK, example_posts)
}