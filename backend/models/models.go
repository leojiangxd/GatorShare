package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Content  string `json:"content"`
	Username string `json:"username"`
}
