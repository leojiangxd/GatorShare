package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gshare.com/platform/models"

	"gshare.com/platform/docs"
)

// @title GatorShare API
// @version 1.0
// @description Backend APIs for the GatorShare app
// @host localhost:8080
// @BasePath /api/v1

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

	docs.SwaggerInfo.Title = "GatorShare API"

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

		// member routes
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
		v1.PUT("post/:postId", updatePost)
		v1.GET("member/:username/posts", getUserPosts)
		v1.PUT("post/:postId/increment-views", incrementPostViews)

		// comment routes
		v1.GET("comment/:postId/", getComments)
		v1.GET("comment/:postId/:commentId", getCommentById)
		v1.POST("comment/:postId", createComment)
		v1.DELETE("comment/:postId/:commentId", deleteComment)

	}

	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}

func index(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "index Called"})

}

// GetMembers godoc
//
//	@Summary		Lists the first 10 members
//	@Description	This API gets the first 10 Member entities from the database
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} string "Success"
//	@Failure 		400 {object} string "Bad Request"
//	@Router			/member [get]
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

// GetMemberByUsername godoc
//
//	@Summary		Gets a member's info by their username
//	@Description	This API fetches a Member entity by their unique username
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200	{object} string "Success"
//	@Failure 		400 {object} string "Bad Request"
//	@Router			/member/{username} [get]
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

// Register godoc
//
//	@Summary		Registers a new member
//	@Description	This API is used to add a Member entity to the database
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member	body		models.Member	true	"New member"
//	@Success		201	{object} string "Created"
//	@Failure 		400 {object} string "Bad Request"
//	@Router			/register [post]
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

	// Check if the email already exists
	if err := db.Select("*").Where("email = ?", newMember.Email).First(&models.Member{}).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
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
	c.JSON(http.StatusCreated, gin.H{"message": "Successfully registered user: " + newMember.Username})
}

// Login godoc
//
//		@Summary		Logs in an existing member
//		@Description	This API is used to login a member by using the stored credentials in the database
//		@Tags			member
//		@Accept			json
//		@Produce		json
//		@Param			member	body	models.Member	true	"Member username and password"
//		@Success		200	{object} string "Success"
//	 	@Failure 		400 {object} string "Bad Request"
//		@Router			/login [post]
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

// Logout godoc
//
//	@Summary		Logs out a currently logged in member
//	@Description	This API clears the tokens in the cookies and database for the logged in Member
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member	body		models.Member	true	"Member username"
//	@Success		200	{object} string "Success"
//	@Failure 		401 {object} string "Unauthorized"
//	@Failure 		404 {object} string "Not Found"
//	@Router			/logout [post]
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

// UpdateMember godoc
//
//	@Summary		Updates a members information
//	@Description	This API updates the field values for the logged-in Member
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member	body		models.Member	true	"Updated member info"
//	@Success		200	{object} string "Success"
//	@Failure 		400 {object} string "Bad Request"
//	@Failure 		401 {object} string "Unauthorized"
//	@Failure 		404 {object} string "Not Found"
//	@Router			/member [put]
func updateMember(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	type UpdateRequest struct {
		CurrentPassword string `json:"currentPassword"`
		NewUsername     string `json:"username"`
		NewEmail        string `json:"email"`
		NewPassword     string `json:"newPassword"`
		Bio             string `json:"bio"`
	}

	var updateReq UpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var currentMember models.Member
	if err := db.First(&currentMember, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if updateReq.NewUsername != "" && updateReq.NewUsername != username {
		var existingUser models.Member
		if err := db.First(&existingUser, "username = ?", updateReq.NewUsername).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	if updateReq.NewEmail != "" && updateReq.NewEmail != currentMember.Email {
		var existingUser models.Member
		if err := db.First(&existingUser, "email = ?", updateReq.NewEmail).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	if updateReq.NewPassword != "" {
		if updateReq.CurrentPassword == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Enter your current password"})
			return
		}
		if !checkPasswordHash(updateReq.CurrentPassword, currentMember.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
			return
		}
		hashedNewPassword, _ := hashPassword(updateReq.NewPassword)
		currentMember.Password = hashedNewPassword
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Update username in posts and comments if changed
		if updateReq.NewUsername != "" && updateReq.NewUsername != username {
			// Update posts
			if err := tx.Model(&models.Post{}).Where("author = ?", username).Update("author", updateReq.NewUsername).Error; err != nil {
				return err
			}

			// Update comments
			if err := tx.Model(&models.Comment{}).Where("author = ?", username).Update("author", updateReq.NewUsername).Error; err != nil {
				return err
			}

			// Update the current member's username
			if err := tx.Model(&currentMember).Update("username", updateReq.NewUsername).Error; err != nil {
				return err
			}
			username = updateReq.NewUsername // Update reference
		}

		if updateReq.NewEmail != "" {
			currentMember.Email = updateReq.NewEmail
		}
		if updateReq.Bio != "" {
			currentMember.Bio = updateReq.Bio
		}
		return tx.Save(&currentMember).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	responseData := gin.H{
		"username": currentMember.Username,
		"email":    currentMember.Email,
		"bio":      currentMember.Bio,
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}

// DeleteMember godoc
//
//	@Summary		Deletes a member from the system
//	@Description	This API removes the Member entity from the database for the logged-in Member
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Param			member	body		models.Member	true	"Member username"
//	@Success		200	{object} string "Success"
//	@Failure 		401 {object} string "Unauthorized"
//	@Failure 		404 {object} string "Not Found"
//	@Router			/member [delete]
func deleteMember(c *gin.Context) {
	// Check if user is logged in
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

	// Start a transaction to ensure all updates happen atomically
	err := db.Transaction(func(tx *gorm.DB) error {
		// Update posts to mark author as [deleted]
		if err := tx.Model(&models.Post{}).Where("author = ?", username).Update("author", "[deleted]").Error; err != nil {
			return err
		}

		// Update comments to mark author as [deleted]
		if err := tx.Model(&models.Comment{}).Where("author = ?", username).Update("author", "[deleted]").Error; err != nil {
			return err
		}

		// Delete the member
		if err := tx.Delete(&member).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

// Options godoc
//
// @Summary			Handles OPTIONS requests
// @Description 	This API handles OPTIONS requests for CORS preflight
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} string "Success"
// @Router 			/member [options]
func options(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "options Called"})

}

// GetCurrentUser godoc
//
// @Summary 		Gets the current logged-in member
// @Description 	This API returns the username of the currently logged-in Member
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} string "Success"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		403 {object} string "Forbidden"
// @Router 			/current-user [get]
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

// GetPosts godoc
//
// @Summary 		Retrieves all posts
// @Description 	This API fetches all posts from the database ordered by creation date
// @Tags 			post
// @Accept 			json
// @Produce 		json
// @Success 		200 {array} models.Post
// @Failure 		400 {object} string "Bad Request"
// @Router 			/post [get]
func getPosts(c *gin.Context) {

	type PostQuery struct {
		Column    string `json:"column"`
		Order     string `json:"order"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		SearchKey string `json:"search_key"`
	}

	//Start by reading in the sorting column and direction
	var postQuery PostQuery
	if err := c.ShouldBindJSON(&postQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//If no limit or offset was passed in, set to -1 for the SQL command
	if postQuery.Limit == 0 {
		postQuery.Limit = -1
	}
	if postQuery.Offset == 0 {
		postQuery.Offset = -1
	}

	//Format the order for sorting
	var order string
	if postQuery.Order != "" {
		order = postQuery.Column + " " + postQuery.Order
	} else {
		order = postQuery.Column
	}

	var posts []models.Post

	// Fetch posts ordered by the passed in column, with slices specified
	result := db.Where("title LIKE ?", "%"+postQuery.SearchKey+"%").Or("author LIKE ?", "%"+postQuery.SearchKey+"%").Or("content LIKE ?", "%"+postQuery.SearchKey+"%").Order(order).Limit(postQuery.Limit).Offset(postQuery.Offset).Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	//Get the count
	var count int64
	db.Model(&models.Post{}).Where("title LIKE ?", "%"+postQuery.SearchKey+"%").Or("author LIKE ?", "%"+postQuery.SearchKey+"%").Or("content LIKE ?", "%"+postQuery.SearchKey+"%").Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count, "data": posts})
}

// GetPostById godoc
//
// @Summary 		Retrieves a specific post by ID
// @Description 	This API fetches a post and its comments by the post ID
// @Tags 			post
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Success 		200 {object} models.Post
// @Failure 		400 {object} string "Bad Request"
// @Router 			/post/{postId} [get]
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

// CreatePost godoc
//
// @Summary 	Creates a new post
// @Description This API creates a new post for the logged-in member
// @Tags 		post
// @Accept 		json
// @Produce 	json
// @Param 		post body models.Post true "New post"
// @Success 	200 {object} models.Post
// @Failure 	400 {object} string "Bad Request"
// @Failure 	401 {object} string "Unauthorized"
// @Router 		/post [post]
func createPost(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the username from the context
	username := getUsername(c)
	post.Author = username

	if post.Title == "" || post.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content is required"})
		return
	}

	if post.Images == nil {
		post.Images = models.StringArray{}
	}
	post.PostId = uuid.New().String()
	post.Comments = []models.Comment{}

	result := db.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "data": post})
}

// DeletePost godoc
//
// @Summary 	Deletes a post
// @Description This API deletes a post by its ID if the logged-in member is the author
// @Tags 		post
// @Accept 		json
// @Produce 	json
// @Param 		postId path string true "Post ID"
// @Success 	200 {object} string
// @Failure 	400 {object} string "Bad Request"
// @Failure 	401 {object} string "Unauthorized"
// @Failure 	403 {object} string "Forbidden"
// @Router 		/post/{postId} [delete]
func deletePost(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post
	if err := db.First(&post, "post_id = ?", c.Param("postId")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the post belongs to the logged-in user
	username := getUsername(c)
	if post.Author != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	db.Delete(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// UpdatePost godoc
//
// @Summary 	Updates a post
// @Description This API updates a post's content if the logged-in member is the author
// @Tags 		post
// @Accept 		json
// @Produce 	json
// @Param 		postId path string true "Post ID"
// @Param 		post body models.Post true "Updated post"
// @Success 	200 {object} models.Post
// @Failure 	400 {object} string "Bad Request"
// @Failure 	401 {object} string "Unauthorized"
// @Failure 	403 {object} string "Forbidden"
// @Router 		/post/{postId} [put]
func updatePost(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var post models.Post
	if err := db.First(&post, "post_id = ?", c.Param("postId")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if the post belongs to the logged-in user
	username := getUsername(c)
	if post.Author != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&post).Where("post_id = ?", c.Param("postId")).Updates(models.Post{Title: post.Title, Content: post.Content})
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "data": post})
}

// GetUserPosts godoc
//
// @Summary 	Retrieves posts by a specific user
// @Description This API fetches all posts created by a specific member
// @Tags 		post
// @Accept 		json
// @Produce 	json
// @Param		member	body		models.Member	true	"Member username"
// @Success 	200 {array} models.Post
// @Failure 	500 {object} string "Internal Server Error"
// @Router 		/member/{username}/posts [get]
func getUserPosts(c *gin.Context) {
	username := c.Param("username")

	var posts []models.Post
	result := db.Where("author = ?", username).Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": posts})
}

// IncrementPostViews godoc
//
// @Summary 	Increments the view count of a post
// @Description This API increments the view count of a specific post by its ID
// @Tags 		post
// @Accept 		json
// @Produce 	json
// @Param 		postId path string true "Post ID"
// @Success 	200 {object} string
// @Failure 	400 {object} string "Bad Request"
// @Router 		/post/{postId}/increment-views [put]
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
