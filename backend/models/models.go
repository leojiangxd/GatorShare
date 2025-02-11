package models

import (
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	MemberId  int    `json:"member_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
