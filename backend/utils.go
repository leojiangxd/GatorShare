package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/bcrypt"

	"time"

	"github.com/google/uuid"
	"gshare.com/platform/models"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func sendAutoNotification(recipient, title, content string) error {
	noti := models.Notification{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		Username:  recipient,
		Title:     title,
		Content:   content,
		Read:      false,
	}
	return db.Create(&noti).Error
}
