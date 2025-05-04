package domain

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         int       `json:"id"`
	UUID       string    `json:"uuid" gorm:"unique;not null"`
	Username   string    `json:"username" gorm:"unique;not null"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Fullname   string    `json:"fullname"`
	Password   string    `json:"-" gorm:"not null"`
	Token      string    `json:"token,omitempty" gorm:"unique"`
	IsVerified bool      `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`
}

// NewUser creates a new User entity
func NewUser(username, email, password string) (*User, error) {
	uuid := uuid.New().String()
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return nil, err
	}

	user := &User{
		UUID:       uuid,
		Username:   username,
		Email:      email,
		Password:   password,
		Token:      signupToken,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	return user, nil
}

// GenerateSignupToken generates a random signup token for the user
func GenerateSignupToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate signup token: %v", err)
	}
	return fmt.Sprintf("%x", token), nil
}
