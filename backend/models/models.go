package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MemberId     int    `json:"member_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"hashed_password"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}
