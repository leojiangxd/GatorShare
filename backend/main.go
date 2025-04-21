package main

import (
	"fmt"
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
	db.AutoMigrate(&models.Member{}, &models.Post{}, &models.Comment{}, &models.Notification{})
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

		v1.GET("member/:username/liked-posts", getUserLikedPosts)
		v1.GET("member/:username/disliked-posts", getUserDislikedPosts)
		v1.GET("member/:username/liked-comments", getUserLikedComments)
		v1.GET("member/:username/disliked-comments", getUserDislikedComments)

		v1.POST("member/:username/follow", followMember)
		v1.DELETE("member/:username/follow", unfollowMember)
		v1.GET("member/:username/followers", getFollowers)
		v1.GET("member/:username/following", getFollowing)

		// post routes
		v1.GET("post", getPosts)
		v1.GET("post/:postId", getPostById)
		v1.POST("post", createPost)
		v1.DELETE("post/:postId", deletePost)
		v1.PUT("post/:postId", updatePost)
		v1.GET("member/:username/posts", getUserPosts)
		v1.PUT("post/:postId/increment-views", incrementPostViews)
		v1.PUT("post/:postId/like-dislike", likeOrDislikePost)

		// comment routes
		v1.GET("comment/:postId/", getComments)
		v1.GET("comment/:postId/:commentId", getCommentById)
		v1.POST("comment/:postId", createComment)
		v1.PUT("comment/:postId/:commentId", updateComment)
		v1.DELETE("comment/:postId/:commentId", deleteComment)
		v1.PUT("comment/:postId/:commentId/like-dislike", likeOrDislikeComment)

		// notification routes
		v1.GET("notification", getNotifications)
		v1.GET("notification/:id", getNotificationById)
		v1.POST("notification", sendNotification)
		v1.DELETE("notification/:id", deleteNotification)
		v1.PUT("notification/:id", updateNotification)
		v1.PUT("notification", updateNotifications)

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
//	@Summary		Gets a list of members
//	@Description	Gets a slice of members using the limit and offset parameters, sorts based on the column and order (desc or asc) parameters, and filters based off the search_key parameter
//	@Tags			member
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} string "Success"
//	@Failure 		400 {object} string "Bad Request"
//	@Router			/member [get]
func getMembers(c *gin.Context) {

	//Start by reading in the sorting column and direction
	var memberQuery models.SearchQuery
	if err := c.ShouldBindQuery(&memberQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//If no limit or offset was passed in, set to -1 for the SQL command
	if memberQuery.Limit == 0 {
		memberQuery.Limit = -1
	}
	if memberQuery.Offset == 0 {
		memberQuery.Offset = -1
	}

	//Format the order for sorting
	var order string
	if memberQuery.Order != "" {
		order = memberQuery.Column + " " + memberQuery.Order
	} else {
		order = memberQuery.Column
	}

	var members []models.Member

	// Fetch posts ordered by the passed in column, with slices specified
	result := db.Where("username LIKE ?", "%"+memberQuery.SearchKey+"%").Or("email LIKE ?", "%"+memberQuery.SearchKey+"%").Or("bio LIKE ?", "%"+memberQuery.SearchKey+"%").Order(order).Limit(memberQuery.Limit).Offset(memberQuery.Offset).Find(&members)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	//Get the count
	var count int64
	db.Model(&models.Member{}).Where("username LIKE ?", "%"+memberQuery.SearchKey+"%").Or("email LIKE ?", "%"+memberQuery.SearchKey+"%").Or("bio LIKE ?", "%"+memberQuery.SearchKey+"%").Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count, "data": members})
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

// GetUserLikedPosts godoc
//
// @Summary 		Retrieves posts liked by a specific user
// @Description 	This API fetches all posts that have been liked by a specific user, identified by their username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {array} models.Post
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "User not found"
// @Router 			/member/{username}/liked-posts [get]
func getUserLikedPosts(c *gin.Context) {
	username := c.Param("username")

	var member models.Member
	if err := db.Preload("LikedPosts").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member.LikedPosts})
}

// GetUserDislikedPosts godoc
//
// @Summary 		Retrieves posts disliked by a specific user
// @Description		This API fetches all posts that have been disliked by a specific user, identified by their username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {array} models.Post
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "User not found"
// @Router 			/member/{username}/disliked-posts [get]
func getUserDislikedPosts(c *gin.Context) {
	username := c.Param("username")

	var member models.Member
	if err := db.Preload("DislikedPosts").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member.DislikedPosts})
}

// GetPosts godoc
//
// @Summary 		Retrieves posts
// @Description 	Gets a slice of posts using the limit and offset parameters, sorts based on the column and order (desc or asc) parameters, and filters based off the search_key parameter
// @Tags 			post
// @Accept 			json
// @Produce 		json
// @Success 		200 {array} models.Post
// @Failure 		400 {object} string "Bad Request"
// @Router 			/post [get]
func getPosts(c *gin.Context) {

	//Start by reading in the sorting column and direction
	var postQuery models.SearchQuery
	if err := c.ShouldBindQuery(&postQuery); err != nil {
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
	if postQuery.Column == "comments" {
		result := db.Preload("Comments").
			Where("title LIKE ?", "%"+postQuery.SearchKey+"%").
			Or("author LIKE ?", "%"+postQuery.SearchKey+"%").
			Or("content LIKE ?", "%"+postQuery.SearchKey+"%").
			Order(fmt.Sprintf("(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.post_id) %s", postQuery.Order)).
			Limit(postQuery.Limit).
			Offset(postQuery.Offset).
			Find(&posts)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}
	} else {
		result := db.Where("title LIKE ?", "%"+postQuery.SearchKey+"%").
			Or("author LIKE ?", "%"+postQuery.SearchKey+"%").
			Or("content LIKE ?", "%"+postQuery.SearchKey+"%").
			Order(order).
			Limit(postQuery.Limit).
			Offset(postQuery.Offset).
			Find(&posts)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
			return
		}
	}

	//Get the count
	var count int64
	db.Model(&models.Post{}).
		Where("title LIKE ?", "%"+postQuery.SearchKey+"%").
		Or("author LIKE ?", "%"+postQuery.SearchKey+"%").
		Or("content LIKE ?", "%"+postQuery.SearchKey+"%").
		Count(&count)

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

	// if post.Title == "" || post.Content == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content is required"})
	// 	return
	// }

	if post.Images == nil {
		post.Images = models.StringArray{}
	}

	// Update post with new title, content, and images
	result := db.Model(&post).Where("post_id = ?", c.Param("postId")).Updates(models.Post{
		Title:   post.Title,
		Content: post.Content,
		Images:  post.Images,
	})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

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

// LikeOrDislikePost godoc
//
// @Summary 		Likes or dislikes a post
// @Description 	This API allows a logged-in user to like or dislike a specific post. The action is specified in the request body as either "like" or "dislike".
// @Tags 			post
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Success 		200 {object} map[string]interface{} "Action applied successfully with updated like/dislike counts"
// @Failure 		400 {object} string "Bad Request or Invalid Action"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		404 {object} string "Post not found"
// @Router 			/post/{postId}/like-dislike [put]
func likeOrDislikePost(c *gin.Context) {
	// Check if the user is authorized
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postId := c.Param("postId")
	var post models.Post

	// Fetch the post by ID
	if err := db.First(&post, "post_id = ?", postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Parse the request body to determine the action (like or dislike)
	var request struct {
		Action string `json:"action"` // "like" or "dislike"
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var member models.Member
	// Fetch the member and preload their liked/disliked posts
	if err := db.Preload("LikedPosts").Preload("DislikedPosts").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	switch request.Action {
	case "like":
		// Check if the user already liked the post
		if containsPost(member.LikedPosts, post.PostId) {
			db.Model(&member).Association("LikedPosts").Delete(&post)
			post.Likes--
		} else {
			// Remove from disliked posts if it exists there
			if containsPost(member.DislikedPosts, post.PostId) {
				db.Model(&member).Association("DislikedPosts").Delete(&post)
				post.Dislikes--
			}
			// Add to liked posts and increment likes counter
			db.Model(&member).Association("LikedPosts").Append(&post)
			post.Likes++
		}
		if post.Author != username {
			title := "Your post was liked!"
			content := fmt.Sprintf("%s liked your post: %s", username, post.Title)
			sendAutoNotification(post.Author, title, content)
		}

	case "dislike":
		// Check if the user already disliked the post
		if containsPost(member.DislikedPosts, post.PostId) {
			db.Model(&member).Association("DislikedPosts").Delete(&post)
			post.Dislikes--
		} else {
			// Remove from liked posts if it exists there
			if containsPost(member.LikedPosts, post.PostId) {
				db.Model(&member).Association("LikedPosts").Delete(&post)
				post.Likes--
			}
			// Add to disliked posts and increment dislikes counter
			db.Model(&member).Association("DislikedPosts").Append(&post)
			post.Dislikes++
		}
		if post.Author != username {
			title := "Your post was disliked!"
			content := fmt.Sprintf("%s disliked your post: %s", username, post.Title)
			sendAutoNotification(post.Author, title, content)
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	// Save the updated post to update likes/dislikes counters in the database
	if err := db.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Action applied successfully",
		"likes":    post.Likes,
		"dislikes": post.Dislikes,
	})
}

// Helper function to check if a post is in a list of posts
func containsPost(posts []*models.Post, postId string) bool {
	for _, p := range posts {
		if p.PostId == postId {
			return true
		}
	}
	return false
}

// GetComments godoc
//
// @Summary 		Retrieves all comments for a specific post
// @Description 	This API fetches all comments associated with a specific post, ordered by creation date (latest first).
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Success 		200 {array} models.Comment "List of comments"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "Post not found"
// @Router 			/comment/{postId}/ [get]
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

// GetCommentById godoc
//
// @Summary 		Retrieves a specific comment by its ID
// @Description 	This API fetches a comment associated with a specific post by its unique comment ID.
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Param 			commentId path string true "Comment ID"
// @Success 		200 {object} models.Comment "Comment details"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "Comment not found"
// @Router 			/comment/{postId}/{commentId} [get]
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

// CreateComment godoc
//
// @Summary 		Creates a new comment on a post
// @Description 	This API allows a logged-in user to add a new comment to a specific post.
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Param 			comment body models.Comment true "New comment data"
// @Success 		200 {object} models.Comment "Created comment details"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		404 {object} string "Post not found"
// @Router 			/comment/{postId} [post]
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
		if post.Author != username {
			title := "New comment on your post!"
			content := fmt.Sprintf("%s commented: %s", username, newComment.Content)
			sendAutoNotification(post.Author, title, content)
		}
	}
}

// UpdateComment godoc
//
// @Summary 		Updates an existing comment
// @Description 	This API allows the author of a comment to update its content.
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Param 			commentId path string true "Comment ID"
// @Success 		200 {object} models.Comment "Updated comment details"
// @Failure 		400 {object} string "Bad Request or Empty Content"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		403 {object} string "Forbidden - Only the author can update their own comments"
// @Failure 		404 {object} string "Comment not found"
// @Router 			/comment/{postId}/{commentId} [put]
func updateComment(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postId := c.Param("postId")
	commentId := c.Param("commentId")

	var comment models.Comment
	if err := db.First(&comment, "comment_id = ? AND post_id = ?", commentId, postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if the comment belongs to the logged-in user
	username := getUsername(c)
	if comment.Author != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own comments"})
		return
	}

	var updateData struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment content cannot be empty"})
		return
	}

	comment.Content = updateData.Content
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"data":    comment,
	})
}

// DeleteComment godoc
//
// @Summary 		Deletes an existing comment
// @Description 	This API allows the author of a comment to delete it from a specific post.
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Param 			commentId path string true "Comment ID"
// @Success 		200 {object} string "Comment deleted successfully"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		403 {object} string "Forbidden - Only the author can delete their own comments"
// @Failure 		404 {object} string "Comment not found"
// @Router 			/comment/{postId}/{commentId} [delete]
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

// LikeOrDislikeComment godoc
//
// @Summary 		Like or dislike action on a comment
// @Description 	This API allows a logged-in user to like or dislike a specific comment. The action is specified in the request body as either "like" or "dislike".
// @Tags 			comment
// @Accept 			json
// @Produce 		json
// @Param 			postId path string true "Post ID"
// @Param 			commentId path string true "Comment ID"
// @Success 		200 {object} map[string]interface{} "Action applied successfully with updated like/dislike counts"
// @Failure 		400 {object} string "Invalid action"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/comment/{postId}/{commentId}/like-dislike [put]
func likeOrDislikeComment(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := getUsername(c)
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postId := c.Param("postId")
	commentId := c.Param("commentId")
	var comment models.Comment

	if err := db.First(&comment, "comment_id = ? AND post_id = ?", commentId, postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	var request struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var member models.Member
	if err := db.Preload("LikedComments").Preload("DislikedComments").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	switch request.Action {
	case "like":
		if containsComment(member.LikedComments, comment.CommentId) {
			db.Model(&member).Association("LikedComments").Delete(&comment)
			comment.Likes--
		} else {
			if containsComment(member.DislikedComments, comment.CommentId) {
				db.Model(&member).Association("DislikedComments").Delete(&comment)
				comment.Dislikes--
			}
			db.Model(&member).Association("LikedComments").Append(&comment)
			comment.Likes++
		}
		if comment.Author != username {
			title := "Your comment was liked!"
			content := fmt.Sprintf("%s liked your comment: %s", username, comment.Content)
			sendAutoNotification(comment.Author, title, content)
		}

	case "dislike":
		if containsComment(member.DislikedComments, comment.CommentId) {
			db.Model(&member).Association("DislikedComments").Delete(&comment)
			comment.Dislikes--
		} else {
			if containsComment(member.LikedComments, comment.CommentId) {
				db.Model(&member).Association("LikedComments").Delete(&comment)
				comment.Likes--
			}
			db.Model(&member).Association("DislikedComments").Append(&comment)
			comment.Dislikes++
		}
		if comment.Author != username {
			title := "Your comment was disliked!"
			content := fmt.Sprintf("%s disliked your comment: %s", username, comment.Content)
			sendAutoNotification(comment.Author, title, content)
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}

	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Action applied successfully",
		"likes":    comment.Likes,
		"dislikes": comment.Dislikes,
	})
}

// Helper function to check if a comment is in a list of comments
func containsComment(comments []*models.Comment, commentId string) bool {
	for _, c := range comments {
		if c.CommentId == commentId {
			return true
		}
	}
	return false
}

// GetUserLikedComments godoc
//
// @Summary 		Retrieves comments liked by a specific user
// @Description 	This API fetches all comments that have been liked by a specific user, identified by their username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {array} models.Comment
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "User not found"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/member/{username}/liked-comments [get]
func getUserLikedComments(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := c.Param("username")

	var member models.Member
	if err := db.Preload("LikedComments").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member.LikedComments})
}

// GetUserDislikedComments godoc
//
// @Summary 		Retrieves comments disliked by a specific user
// @Description 	This API fetches all comments that have been disliked by a specific user, identified by their username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {array} models.Comment
// @Failure 		400 {object} string "Bad Request"
// @Failure 		404 {object} string "User not found"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/member/{username}/disliked-comments [get]
func getUserDislikedComments(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	username := c.Param("username")

	var member models.Member
	if err := db.Preload("DislikedComments").First(&member, "username = ?", username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": member.DislikedComments})
}

// GetNotifications godoc
//
// @Summary 		Gets notifications for the current user
// @Description 	Fetches a slice of notifications for the logged-in user. Supports optional query parameters for sorting, limit, and offset.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} map[string]interface{} "Notifications list and count"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification [get]
func getNotifications(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	//Start by reading in the sorting column and direction
	var notiQuery models.SearchQuery
	if err := c.ShouldBindQuery(&notiQuery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//If no limit or offset was passed in, set to -1 for the SQL command
	if notiQuery.Limit == 0 {
		notiQuery.Limit = -1
	}
	if notiQuery.Offset == 0 {
		notiQuery.Offset = -1
	}

	//Format the order for sorting
	var order string
	if notiQuery.Order != "" {
		order = notiQuery.Column + " " + notiQuery.Order
	} else {
		order = notiQuery.Column
	}

	var notis []models.Notification

	// Fetch posts ordered by the passed in column, with slices specified
	result := db.Where("username = ?", getUsername(c)).Order(order).Limit(notiQuery.Limit).Offset(notiQuery.Offset).Find(&notis)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	//Get the count
	var count int64
	db.Model(&models.Notification{}).Where("username = ?", getUsername(c)).Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count, "data": notis})
}

// GetNotificationById godoc
//
// @Summary 		Gets a notification by its ID
// @Description 	Fetches a single notification by its unique ID.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Notification ID"
// @Success 		200 {object} models.Notification
// @Failure 		400 {object} string "No Records Found"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification/{id} [get]
func getNotificationById(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var noti models.Notification

	result := db.First(&noti, "id = ?", id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": noti})
	}
}

// SendNotification godoc
//
// @Summary 		Sends a notification to a user
// @Description 	Sends a notification with a title and content to a specified user.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Param 			notification body models.Notification true "Notification data"
// @Success 		201 {object} string "Notification sent"
// @Failure 		400 {object} string "Notification title, content, and username (recipient) required"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification [post]
func sendNotification(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var noti models.Notification
	if err := c.ShouldBindJSON(&noti); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if noti.Title == "" || noti.Content == "" || noti.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification title, content, and username (recipient) required"})
		return
	}

	noti.Id = uuid.New().String()
	noti.Read = false

	db.Create(&noti)

	c.JSON(http.StatusCreated, gin.H{"message": "Notification sent"})
}

// DeleteNotification godoc
//
// @Summary 		Deletes a notification
// @Description 	Deletes a notification by its unique ID. Only the recipient can delete their own notifications.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Notification ID"
// @Success 		200 {object} string "Notification deleted successfully"
// @Failure 		400 {object} string "Notification not found or not owned by user"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification/{id} [delete]
func deleteNotification(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var noti models.Notification
	if err := db.First(&noti, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notification not found"})
		return
	}

	if noti.Username != getUsername(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only delete your own notifications"})
		return
	}

	db.Delete(&noti)
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}

// UpdateNotification godoc
//
// @Summary 		Updates a notification's read status
// @Description 	Updates the read status of a notification by its unique ID. Only the recipient can update their own notifications.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Notification ID"
// @Success 		200 {object} string "Notification updated successfully"
// @Failure 		404 {object} string "Notification not found"
// @Failure 		400 {object} string "You can only update your own notifications"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification/{id} [put]
func updateNotification(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")

	var noti models.Notification
	if err := db.First(&noti, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	if noti.Username != getUsername(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only update your own notifications"})
		return
	}

	var updateReq struct {
		Read bool `json:"read"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&models.Notification{}).Where("id = ?", id).Update("read", updateReq.Read)
	c.JSON(http.StatusOK, gin.H{"message": "Notification updated successfully"})
}

// UpdateNotifications godoc
//
// @Summary 		Bulk update notifications' read status
// @Description 	Sets all notifications for the logged-in user to read/unread.
// @Tags 			notification
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} string "Notifications updated successfully"
// @Failure 		400 {object} string "Bad Request"
// @Failure 		401 {object} string "Unauthorized"
// @Router 			/notification [put]
func updateNotifications(c *gin.Context) {

	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updateReq struct {
		Read bool `json:"read"`
	}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&models.Notification{}).Where("username = ?", getUsername(c)).Update("read", updateReq.Read)
	c.JSON(http.StatusOK, gin.H{"message": "Notifications updated successfully"})
}

// FollowMember godoc
//
// @Summary 		Follow a member
// @Description 	Allows the logged-in user to follow another member by username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member to follow"
// @Success 		200 {object} string "Followed successfully"
// @Failure 		400 {object} string "Cannot follow yourself or bad request"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		404 {object} string "User not found"
// @Failure 		500 {object} string "Failed to follow"
// @Router 			/member/{username}/follow [post]
func followMember(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	follower := getUsername(c)
	followee := c.Param("username")
	if follower == followee {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
		return
	}

	var followerMember, followeeMember models.Member
	if db.First(&followerMember, "username = ?", follower).Error != nil ||
		db.First(&followeeMember, "username = ?", followee).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Add follow relationship
	if err := db.Model(&followerMember).Association("Following").Append(&followeeMember); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow"})
		return
	}

	// Optionally, send notification
	if follower != followee {
		sendAutoNotification(followee, "New follower!", follower+" started following you.")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

// UnfollowMember godoc
//
// @Summary 		Unfollow a member
// @Description 	Allows the logged-in user to unfollow another member by username.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member to unfollow"
// @Success 		200 {object} string "Unfollowed successfully"
// @Failure 		400 {object} string "Cannot unfollow yourself or bad request"
// @Failure 		401 {object} string "Unauthorized"
// @Failure 		404 {object} string "User not found"
// @Failure 		500 {object} string "Failed to unfollow"
// @Router 			/member/{username}/follow [delete]
func unfollowMember(c *gin.Context) {
	if err := Authorize(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	follower := getUsername(c)
	followee := c.Param("username")
	if follower == followee {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot unfollow yourself"})
		return
	}

	var followerMember, followeeMember models.Member
	if db.First(&followerMember, "username = ?", follower).Error != nil ||
		db.First(&followeeMember, "username = ?", followee).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Remove follow relationship
	if err := db.Model(&followerMember).Association("Following").Delete(&followeeMember); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

// GetFollowers godoc
//
// @Summary 		Get a member's followers
// @Description 	Retrieves a list of members who follow the specified user.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {object} map[string]interface{} "List of followers"
// @Failure 		404 {object} string "User not found"
// @Router 			/member/{username}/followers [get]
func getFollowers(c *gin.Context) {
	username := c.Param("username")
	var member models.Member
	if db.Preload("Followers").First(&member, "username = ?", username).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": member.Followers})
}

// GetFollowing godoc
//
// @Summary 		Get members a user is following
// @Description 	Retrieves a list of members that the specified user is following.
// @Tags 			member
// @Accept 			json
// @Produce 		json
// @Param 			username path string true "Username of the member"
// @Success 		200 {object} map[string]interface{} "List of following"
// @Failure 		404 {object} string "User not found"
// @Router 			/member/{username}/following [get]
func getFollowing(c *gin.Context) {
	username := c.Param("username")
	var member models.Member
	if db.Preload("Following").First(&member, "username = ?", username).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": member.Following})
}
