package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Email        string `json:"email"`
	Username     string `json:"username" gorm:"primaryKey"`
	Password     string `json:"password"`
	Bio          string `json:"bio"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}

type StringArray []string

func (sa *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Failed to unmarshal StringArray value")
	}

	return json.Unmarshal(bytes, &sa)
}

func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
}

type Post struct {
	gorm.Model
	PostId   string      `json:"post_id" gorm:"primaryKey"`
	Author   string      `json:"author"`
	Title    string      `json:"title"`
	Content  string      `json:"content"`
	Likes    int         `json:"likes"`
	Dislikes int         `json:"dislikes"`
	Views    int         `json:"views"`
	Comments []Comment   `json:"comments" gorm:"foreignKey:PostID;references:PostId"`
	Images   StringArray `json:"images" gorm:"type:text"`
}

type Comment struct {
	gorm.Model
	CommentId string `json:"comment_id" gorm:"primaryKey"`
	PostID    string `json:"post_id" gorm:"index"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
}
