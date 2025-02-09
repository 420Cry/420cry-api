package controllers

import (
	models "cry-api/app/api/models/users"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateSignupToken generates a random signup token
func GenerateSignupToken() (string, error) {
	token := make([]byte, 32) // 32-byte secure token
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate signup token: %v", err)
	}
	return fmt.Sprintf("%x", token), nil
}

// checkUserExistence checks if a user with the given username or email exists
func checkUserExistence(db *gorm.DB, username, email string) (*models.User, error) {
	var existingUser models.User
	err := db.Where("username = ?", username).Or("email = ?", email).First(&existingUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &existingUser, nil
}

// CreateUser creates a new user in the database with a generated UUID and signup token
func CreateUser(db *gorm.DB, newUser models.User) (models.User, error) {
	// Check if the username or email already exists
	existingUser, err := checkUserExistence(db, newUser.Username, newUser.Email)
	if err != nil {
		return models.User{}, err
	}

	if existingUser != nil {
		// Determine which field is already taken
		if existingUser.Username == newUser.Username {
			return models.User{}, fmt.Errorf("username is already taken")
		}
		if existingUser.Email == newUser.Email {
			return models.User{}, fmt.Errorf("email is already taken")
		}
	}

	// Generate a unique UUID
	newUser.UUID = uuid.New().String()

	// Generate a unique signup token
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return models.User{}, err
	}
	newUser.Token = signupToken

	// Create the user in the DB
	if err := db.Create(&newUser).Error; err != nil {
		// Handle unique constraint violation for MySQL (error code 1062)
		if strings.Contains(err.Error(), "1062") {
			return models.User{}, fmt.Errorf("duplicate entry detected")
		}
		return models.User{}, fmt.Errorf("internal server error: %v", err)
	}

	// Successfully created user
	return newUser, nil
}
