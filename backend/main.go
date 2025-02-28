package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
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
	db.AutoMigrate(&models.Member{}, &models.Post{}, &models.Comment{})
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "X-CSRF-Token", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", index)
		
		// user routes
		v1.GET("member", getMembers)
		v1.GET("member/:username", getMemberByUsername)
		v1.POST("register", register)
		v1.PUT("member", updateMember)
		v1.DELETE("member", deleteMember)
		v1.POST("login", login)
		v1.POST("logout", logout)
		v1.OPTIONS("member", options)

		v1.GET("current-user", getCurrentUser)
		
		// post routes
		v1.GET("post", getPosts)
		v1.GET("post/:postId", getPostById)
		v1.POST("post", createPost)
		v1.DELETE("post/:postId", deletePost)
		v1.PUT("post/:postId/increment-views", incrementPostViews)

		// comment routes
		v1.GET("comment/:postId/", getComments)
		v1.GET("comment/:postId/:commentId", getCommentById)
		v1.POST("comment/:postId", createComment)
		v1.DELETE("comment/:postId/:commentId", deleteComment)

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
	}
	
	//Generate tokens and pass to user cookies
	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	c.SetCookie("session_token", sessionToken, 3600, "/", "localhost", false, true)
	c.SetCookie("csrf_token", csrfToken, 3600, "/", "localhost", false, false)

	//Store tokens in the database
	db.Model(&newMember).Where("username = ?", newMember.Username).Updates(map[string]any{"session_token": sessionToken, "csrf_token": csrfToken})
	c.JSON(http.StatusOK, gin.H{"message": "Register successful"})
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

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

func getCurrentUser(c *gin.Context) {
	// Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized 2"})
		return
	}
	
	// Return the username
	c.JSON(http.StatusOK, gin.H{"username": username})
}

func getPosts(c *gin.Context) {
	var posts []models.Post

	// Fetch posts ordered by createdAt descending (latest first)
	result := db.Preload("Comments").Order("created_at desc").Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	if len(posts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Posts Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

func getPostById(c *gin.Context) {
	postId := c.Param("postId")
	var post models.Post

	// Fetch post and preload its comments ordered by createdAt descending
	result := db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&post, "post_id = ?", postId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": post})
}


func createPost(c *gin.Context) {
	//Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var newPost models.Post

	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newPost.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	newPost.Author = username
	newPost.Likes = 0
	newPost.Dislikes = 0
	newPost.Views = 0
	newPost.PostId = uuid.New().String()
	newPost.Comments = []models.Comment{}

	result := db.Create(&newPost)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "data": newPost})
	}
}

func deletePost(c *gin.Context) {

	//Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// find post by id
	postId := c.Param("postId")
	var post models.Post
	result := db.First(&post, "post_id = ?", postId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": post})
	}

	if post.Author != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can only delete your own posts"})
		return
	}

	db.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func incrementPostViews(c *gin.Context) {
	postId := c.Param("postId")

	// Find the post by ID
	var post models.Post
	result := db.First(&post, "post_id = ?", postId)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	}

	// Increment the view count by 1
	post.Views++

	// Save the updated post
	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update view count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "View count incremented", "views": post.Views})
}

func getComments(c *gin.Context) {
	postId := c.Param("postId")
	
	// Find the post first
	var post models.Post
	if err := db.First(&post, "post_id = ?", postId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	}

	// Preload comments for the post ordered by createdAt descending (latest first)
	var comments []models.Comment
	result := db.Order("created_at desc").Where("post_id = ?", postId).Find(&comments)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func getCommentById(c *gin.Context) {
	postId := c.Param("postId")
	commentId := c.Param("commentId")
	
	var comment models.Comment
	result := db.First(&comment, "comment_id = ? AND post_id = ?", commentId, postId)
	
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": comment})
	}
}

func createComment(c *gin.Context) {
	// Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Get the post ID from URL parameter
	postId := c.Param("postId")
	
	// Verify the post exists
	var post models.Post
	if err := db.First(&post, "post_id = ?", postId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Post not found"})
		return
	}
	
	var newComment models.Comment
	
	// Bind the JSON request to the comment struct
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate comment content
	if newComment.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content is required"})
		return
	}
	
	// Set comment fields
	newComment.Author = username
	newComment.PostID = postId
	newComment.CommentId = uuid.New().String()
	newComment.Likes = 0
	newComment.Dislikes = 0
	
	// Create the comment in the database
	result := db.Create(&newComment)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Comment created successfully", "data": newComment})
	}
}

func deleteComment(c *gin.Context) {
	// Check if user is logged in
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	// Get parameters from URL
	postId := c.Param("postId")
	commentId := c.Param("commentId")
	
	// Find the comment
	var comment models.Comment
	result := db.First(&comment, "comment_id = ? AND post_id = ?", commentId, postId)
	
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment not found"})
		return
	}
	
	// Check if the user is the author of the comment
	if comment.Author != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You can only delete your own comments"})
		return
	}
	
	// Delete the comment
	db.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}