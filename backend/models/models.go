package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

	// Relationships
	LikedPosts       []*Post    `gorm:"many2many:member_likes;" json:"liked_posts"`
	DislikedPosts    []*Post    `gorm:"many2many:member_dislikes;" json:"disliked_posts"`
	LikedComments    []*Comment `gorm:"many2many:member_comment_likes;" json:"liked_comments"`
	DislikedComments []*Comment `gorm:"many2many:member_comment_dislikes;" json:"disliked_comments"`
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
	PostId    string `json:"post_id" gorm:"primaryKey"`
	CreatedAt time.Time
	Author    string      `json:"author"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Likes     int         `json:"likes"`
	Dislikes  int         `json:"dislikes"`
	Views     int         `json:"views"`
	Comments  []Comment   `json:"comments" gorm:"foreignKey:PostID;references:PostId"`
	Images    StringArray `json:"images" gorm:"type:text"`

	// Relationships
	LikedByMembers    []*Member `gorm:"many2many:member_likes;" json:"liked_by_members"`
	DislikedByMembers []*Member `gorm:"many2many:member_dislikes;" json:"disliked_by_members"`
}

type Comment struct {
	CommentId string `json:"comment_id" gorm:"primaryKey"`
	PostID    string `json:"post_id" gorm:"index"`
	CreatedAt time.Time
	Author    string `json:"author"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`

	// Relationships
	LikedByMembers    []*Member `gorm:"many2many:member_comment_likes;" json:"liked_comments"`
	DislikedByMembers []*Member `gorm:"many2many:member_comment_dislikes;" json:"disliked_comments"`
}

type Notification struct {
	Id        string `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Username  string `json:"username"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Read      bool   `json:"read"`
}

type SearchQuery struct {
	Column    string `form:"column"`
	Order     string `form:"order"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
	SearchKey string `form:"search_key"`
}
