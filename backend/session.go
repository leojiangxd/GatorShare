package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gshare.com/platform/models"
)

func Authorize(c *gin.Context) error {
	username := c.Param("username")
	log.Println("Authorize: username =", username)

	// Get the user
	var member models.Member
	result := db.First(&member, "username = ?", username)
	if result.Error != nil {
		log.Println("Authorize error: No such user")
		return errors.New("error: No such user")
	}
	log.Println("Authorize: found user =", member.Username)

	// Check session token
	st, err := c.Cookie("session_token")
	if err != nil {
		log.Println("Authorize error: session_token not found in cookies")
		return errors.New("Invalid session token")
	}
	log.Println("Authorize: session_token from cookie =", st)
	if st == "" || st != member.SessionToken {
		log.Println("Authorize error: session_token does not match")
		return errors.New("Invalid session token")
	}

	// Check the CSRF token from the headers
	csrf, err := c.Cookie("csrf_token")
	if err != nil || csrf != member.CSRFToken || csrf == "" {
		log.Println("Authorize error: Invalid csrf token")
		return errors.New("Invalid csrf token")
	}
	log.Println("Authorize: csrf_token from header =", csrf)
	if csrf != member.CSRFToken {
		log.Println("Authorize error: csrf_token does not match")
		return errors.New("Invalid csrf token")
	}

	// Set the username in the context
	c.Set("username", member.Username)
	log.Println("Authorize: user authorized successfully")
	return nil
}
