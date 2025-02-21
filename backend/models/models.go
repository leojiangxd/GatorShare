package models

import (
	"time"
)

type Member struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	MemberId     int    `json:"member_id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SessionToken string `json:"session_token"`
	CSRFToken    string `json:"csrf_token"`
}
