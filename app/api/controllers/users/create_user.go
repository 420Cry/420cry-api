package controllers

import (
	types "cry-api/app/types/users"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateSignupToken generates a random signup token
func GenerateSignupToken() (string, error) {
	token := make([]byte, 32) // Adjust size as needed
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate signup token: %v", err)
	}
	return fmt.Sprintf("%x", token), nil
}

// CreateUser creates a new user in the database with a generated UUID and signup token
func CreateUser(db *gorm.DB, user types.User) (types.User, error) {
	// Generate a unique UUID
	user.UUID = uuid.New().String()

	// Generate a unique signup token
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return types.User{}, err
	}
	user.SignupToken = signupToken

	// Create the user in the DB
	if err := db.Create(&user).Error; err != nil {
		// Check for duplicate entry error (error code 1062) for MySQL
		if strings.Contains(err.Error(), "1062") {
			// Specific handling for unique constraint violation (username or UUID)
			return types.User{}, fmt.Errorf("duplicate username or UUID")
		}

		// Return any other error
		return types.User{}, fmt.Errorf("internal server error: %v", err)
	}

	// Successfully created user
	return user, nil
}
