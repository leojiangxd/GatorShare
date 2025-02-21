package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gshare.com/platform/models"
)

func Authorize(c *gin.Context) error {

	username := c.Param("username")

	//Get the user
	var member models.Member
	result := db.First(&member, "username = ?", username)
	if result == nil {
		return errors.New("error: No such user")
	}

	//Check session token
	st, err := c.Cookie("session_token")
	if err != nil || st == "" || st != member.SessionToken {
		return errors.New("Invalid session token")
	}

	//Check the CSRF token from the headers
	csrf := c.Request.Header.Get("X-CSRF-Token")
	if csrf != member.CSRFToken || csrf == "" {
		return errors.New("Invalid csrf token")
	}
	return nil
}
