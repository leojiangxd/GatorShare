package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gshare.com/platform/models"
)

func Authorize(c *gin.Context) error {

	authError := errors.New("Unauthorized")

	username := getUsername(c)

	// Get the user
	var member models.Member
	result := db.First(&member, "username = ?", username)
	if result.Error != nil {
		//log.Println("Authorize error: No such user")
		return errors.New("error: No such user")
	}
	//log.Println("Authorize: found user =", member.Username)

	// Check session token
	st, err := c.Cookie("session_token")
	if err != nil {
		//log.Println("Authorize error: session_token not found in cookies")
		return authError
	}
	//log.Println("Authorize: session_token from cookie =", st)
	if st == "" || st != member.SessionToken {
		log.Println("Authorize error: session_token does not match")
		return authError
	}

	// Check the CSRF token from the headers
	csrf := c.Request.Header.Get("X-CSRF-Token")
	if csrf != member.CSRFToken || csrf == "" {
		//log.Println("Authorize error: Invalid csrf token")
		return authError
	}
	//log.Println("Authorize: csrf_token from header =", csrf)
	if csrf != member.CSRFToken {
		//log.Println("Authorize error: csrf_token does not match")
		return authError
	}

	// Set the username in the context
	c.Set("username", member.Username)
	//log.Println("Authorize: user authorized successfully")
	return nil
}

func getUsername(c *gin.Context) string {
	st, err := c.Cookie("session_token")
	if err != nil || st == "" {
		return ""
	}

	var member models.Member
	if db.First(&member, "session_token = ?", st).Error != nil {
		return ""
	}

	return member.Username
}
