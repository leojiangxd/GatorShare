package models

import (
	"time"
)

type Member struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Email        string `json:"email"`
	Username     string `json:"username" gorm:"primaryKey"`		
	Password     string `json:"password"`
	Bio          string `json:"bio"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}

type Post struct {
    PostId       string    `json:"post_id" gorm:"primaryKey"`
    CreatedAt    time.Time
    Author       string    `json:"author"`
    Title        string    `json:"title"`
    Content      string    `json:"content"`
    Likes        int       `json:"likes"`
    Dislikes     int       `json:"dislikes"`
    Views        int       `json:"views"`
    Comments     []Comment `json:"comments" gorm:"foreignKey:PostID;references:PostId"`
}

type Comment struct {
    CommentId    string    `json:"comment_id" gorm:"primaryKey"`
    PostID       string    `json:"post_id" gorm:"index"`
    CreatedAt    time.Time
    Author       string    `json:"author"`
    Content      string    `json:"content"`
    Likes        int       `json:"likes"`
    Dislikes     int       `json:"dislikes"`
}