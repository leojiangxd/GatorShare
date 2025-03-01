package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gshare.com/platform/models"
)

var testSessionToken string
var testCSRFToken string

func SetUpRouter() *gin.Engine {
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
		v1.PUT("post/:postId", updatePost)
		v1.GET("member/:username/posts", getUserPosts)
		v1.PUT("post/:postId/increment-views", incrementPostViews)

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

func TestGetByUsername(t *testing.T) {
	err := connectDatabase()
	checkErr(err)
	r := SetUpRouter()

	req, _ := http.NewRequest("GET", "/api/v1/member/saul", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	user := &models.Member{}
	mockUser := models.Member{
		Email:    "bettercallsaul@test.com",
		Username: "saul",
		Password: "Money123",
		Bio:      "Test Bio",
	}

	log.Println(w.Body.Bytes())
	json.Unmarshal(w.Body.Bytes(), &user)

	assert.Equal(t, mockUser.Email, user.Email)
	assert.Equal(t, mockUser.Username, user.Username)
	assert.True(t, checkPasswordHash(mockUser.Password, user.Password))
	assert.Equal(t, mockUser.Bio, user.Bio)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser(t *testing.T) {
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
