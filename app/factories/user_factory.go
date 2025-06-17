// Package factories provides factory functions for creating and initializing
package factories

import (
	"crypto/rand"
	"fmt"
	"time"

	"cry-api/app/models"
	services "cry-api/app/services/password"

	"github.com/google/uuid"
)

// NewUser creates a new User model instance with a hashed password,
// a unique UUID, a signup token, and an email verification token (OTP).
func NewUser(fullname, username, email, password string) (*models.User, error) {
	// Generate a new UUID for the user
	u := uuid.New().String()

	// Generate a secure signup token (used in activation links)
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return nil, err
	}

	// Generate a short verification token (OTP) for email verification
	verificationToken, err := GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Hash the plaintext password using bcrypt
	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Assemble the User model with all required fields
	user := &models.User{
		UUID:                       u,
		Username:                   username,
		Fullname:                   fullname,
		Email:                      email,
		Password:                   hashedPassword,
		Token:                      &signupToken,      // Signup token for account activation
		VerificationTokens:         verificationToken, // Email verification token (OTP)
		VerificationTokenCreatedAt: time.Now(),        // Timestamp of token creation
		IsVerified:                 false,             // User starts as unverified
		CreatedAt:                  time.Now(),        // Creation timestamp
	}
	return user, nil
}

// GenerateSignupToken creates a secure random 32-byte token returned as a hex string.
// This token is typically embedded in account activation URLs.
func GenerateSignupToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate signup token: %v", err)
	}
	return fmt.Sprintf("%x", token), nil
}

// GenerateVerificationToken creates a 6-character alphanumeric token using uppercase letters and digits.
// This token serves as a short verification code (OTP) typically sent to user emails.
func GenerateVerificationToken() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	b := make([]byte, length)
	randBytes := make([]byte, length)

	// Read random bytes securely
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Map random bytes to allowed charset characters
	for i := range b {
		b[i] = charset[int(randBytes[i])%len(charset)]
	}
	return string(b), nil
}
