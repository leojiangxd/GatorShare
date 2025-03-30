package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gshare.com/platform/models"
)

var testSessionToken string
var testCSRFToken string

func SetUpRouter() *gin.Engine {

	//DEBUG CODE FOR RUNNING INDIVIDUAL UNIT TESTS -------------- DELETE
	//testSessionToken = "B150Bbo7ogNjSvakXtWLHwBRa8FerS1K9SK7GPDhIos="
	//testCSRFToken = "f3L_w5ZcPNLa9vAw8NCIUj5xegjcu0kNfBxguDjzmHA="
	//DEBUG CODE FOR RUNNING INDIVIDUAL UNIT TESTS -------------- DELETE

	r := gin.Default()

	//COMMENT ONE TO CHANGE BETWEEN RUN MODES---------------------------------------------------------------
	gin.SetMode(gin.ReleaseMode)
	//gin.SetMode(gin.DebugMode)
	//COMMENT ONE TO CHANGE BETWEEN RUN MODES---------------------------------------------------------------

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
		v1.GET("member/:username/liked-posts", getUserLikedPosts)

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
		v1.DELETE("comment/:postId/:commentId", deleteComment)

	}
	return r
}

func TestRegister(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	user := models.Member{
		Email:    "bettercallsaul@test.com",
		Username: "saul",
		Password: "Money123",
		Bio:      "Test Bio",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"message":"Successfully registered user: saul"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetMembers(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/member", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := w.Body.String()[8 : len(w.Body.String())-1]

	assert.NotEmpty(t, response)
}

func TestGetMemberByUsername(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	var user models.Member
	mockUser := models.Member{
		Email:    "bettercallsaul@test.com",
		Username: "saul",
		Password: "Money123",
		Bio:      "Test Bio",
	}

	req, _ := http.NewRequest("GET", "/api/v1/member/saul", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := w.Body.String()[8 : len(w.Body.String())-1]
	json.Unmarshal([]byte(response), &user)

	testSessionToken = user.SessionToken
	testCSRFToken = user.CSRFToken

	assert.Equal(t, mockUser.Email, user.Email)
	assert.Equal(t, mockUser.Username, user.Username)
	assert.True(t, checkPasswordHash(mockUser.Password, user.Password))
	assert.Equal(t, mockUser.Bio, user.Bio)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogout(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	user := models.Member{
		Username: "saul",
		Password: "Money123",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/logout", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:     "session_token",
		Value:    testSessionToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: true,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"message":"Logout successful"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	user := models.Member{
		Username: "saul",
		Password: "Money123",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	cookies := w.Result().Cookies()
	testSessionToken = cookies[0].Value[:len(cookies[0].Value)-3] + "="
	testCSRFToken = cookies[1].Value[:len(cookies[1].Value)-3] + "="

	mockResponse := `{"message":"Login successful"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateMember(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	var user models.Member
	mockUser := models.Member{
		Email:    "updated@test.com",
		Username: "saul",
		Password: "Lawyering",
		Bio:      "Updated Bio",
	}

	jsonValue, _ := json.Marshal(mockUser)
	req, _ := http.NewRequest("PUT", "/api/v1/member", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:     "session_token",
		Value:    testSessionToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: true,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	response := w.Body.String()[8 : len(w.Body.String())-1]
	json.Unmarshal([]byte(response), &user)

	assert.Equal(t, mockUser.Email, user.Email)
	assert.Equal(t, mockUser.Username, user.Username)
	assert.True(t, checkPasswordHash(mockUser.Password, user.Password))
	assert.Equal(t, mockUser.Bio, user.Bio)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCurrentUser(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/current-user", nil)
	req.AddCookie(&http.Cookie{
		Name:     "session_token",
		Value:    testSessionToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: true,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"username":"saul"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMember(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	user := models.Member{
		Username: "saul",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("DELETE", "/api/v1/member", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:     "session_token",
		Value:    testSessionToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour),
		Secure:   false,
		HttpOnly: true,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	mockResponse := `{"message":"Member deleted successfully"}`
	responseData, _ := io.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreatePost(t *testing.T) {
	TestRegister(t)
	TestLogin(t)
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	newPost := models.Post{
		Title:   "Saul's create post",
		Content: "To test createPost's functionality",
	}

	jsonValue, _ := json.Marshal(newPost)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post created successfully")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	createdPost := response["data"].(map[string]interface{})
	assert.NotEmpty(t, createdPost["post_id"])
}

func TestGetPosts(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/post", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	posts := response["data"].([]interface{})
	assert.NotEmpty(t, posts)
}

func TestGetPostById(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Creating a new post
	newPost := models.Post{
		Title:   "Saul's get post",
		Content: "To test getPostById functionality",
	}
	jsonValue, _ := json.Marshal(newPost)

	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Extracting post ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	req, _ = http.NewRequest("GET", "/api/v1/post/"+postID, nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	post := response["data"].(map[string]interface{})
	assert.Equal(t, postID, post["post_id"])
}

func TestUpdatePost(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Creating a new post
	newPost := models.Post{
		Title:   "Saul's update post",
		Content: ".....to test updatePost's functionality......",
	}
	jsonValue, _ := json.Marshal(newPost)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Extracting post ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	// Updating the post
	updatedPost := models.Post{
		Content: "To test updatePost's functionality",
	}
	jsonValue, _ = json.Marshal(updatedPost)

	req, _ = http.NewRequest("PUT", "/api/v1/post/"+postID, bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post updated successfully")
}

func TestDeletePost(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Creating a temporary post
	newPost := models.Post{
		Title:   "Saul's deleting post",
		Content: "to test deletePost's functionality",
	}
	jsonValue, _ := json.Marshal(newPost)

	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Extracting the post ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	// Deleting the post
	req, _ = http.NewRequest("DELETE", "/api/v1/post/"+postID, nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post deleted successfully")

}

func TestGetUserPosts(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/member/saul/posts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	posts := response["data"].([]interface{})
	assert.NotEmpty(t, posts)
}

func TestIncrementPostViews(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Creating a new post
	newPost := models.Post{
		Title:   "Saul's increment view post",
		Content: "To test incrementPostViews' functionality",
	}
	jsonValue, _ := json.Marshal(newPost)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))

	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: testSessionToken,
	})
	req.Header.Add("X-CSRF-Token", testCSRFToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Extracting post ID from the response
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	req, _ = http.NewRequest("PUT", "/api/v1/post/"+postID+"/increment-views", nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "View count incremented")

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["views"])

}

func TestLikePost(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Create a new post for testing
	newPost := models.Post{
		Title:   "Test Like Post",
		Content: "Testing like functionality",
	}
	jsonValue, _ := json.Marshal(newPost)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	// Like the post
	likeRequest := map[string]string{"action": "like"}
	jsonValue, _ = json.Marshal(likeRequest)
	req, _ = http.NewRequest("PUT", "/api/v1/post/"+postID+"/like-dislike", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Action applied successfully")
}

func TestGetUserLikedPosts(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/member/saul/liked-posts", nil)
	req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestDislikePost(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	// Create a new post for testing
	newPost := models.Post{
		Title:   "Test Dislike Post",
		Content: "Testing dislike functionality",
	}
	jsonValue, _ := json.Marshal(newPost)
	req, _ := http.NewRequest("POST", "/api/v1/post", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	postID := createResponse["data"].(map[string]interface{})["post_id"].(string)

	// Dislike the post
	dislikeRequest := map[string]string{"action": "dislike"}
	jsonValue, _ = json.Marshal(dislikeRequest)
	req, _ = http.NewRequest("PUT", "/api/v1/post/"+postID+"/like-dislike", bytes.NewBuffer(jsonValue))
	req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
	req.Header.Add("X-CSRF-Token", testCSRFToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Action applied successfully")

	//Deleting user created for testing APIs
	TestDeleteMember(t)
}
