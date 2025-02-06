package api

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	PostID		int		`json:"post_id"`
	UserID		int		`json:"user_id"`
	Title		string	`json:"title"`
	Content		string	`json:"content"`
	Date		string	`json:"date"`
	Likes		int		`json:"likes"`
	Dislikes	int		`json:"dislikes"`
}

// example posts
// var example_posts = []Post{
// 	{PostID: 1, UserID: 1, Title: "Post Title", Content: "Lorem Ipsum", Date: "2025-01-01 00:00:00", Likes: 101, Dislikes: 100},
// 	{PostID: 2, UserID: 2, Title: "Example Post", Content: "Lorem Ipsum", Date: "2025-01-01 00:00:00", Likes: 50, Dislikes: 0},
// }

func GetPosts(c *gin.Context) {
	fmt.Println("Getting posts...")
	// c.IndentedJSON(http.StatusOK, example_posts)

	dbPath := filepath.Join("..", "db", "database.sqlite")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error opening database"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Date, &post.Likes, &post.Dislikes)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through rows"})
	}

	c.IndentedJSON(http.StatusOK, posts)
}