// Package factories provides factory functions for creating and initializing
package factories

import (
	"time"

	"cry-api/app/models"
	PasswordService "cry-api/app/services/password"

	"github.com/google/uuid"
)

// NewUser creates a new instance of models.User with the provided fullname, username, email, and password.
// It generates a UUID, signup token, and verification token for the user, and hashes the provided password.
// Returns the created User object or an error if any step fails.
func NewUser(fullname, username, email, password string) (*models.User, error) {
	u := generateUUID()
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return nil, err
	}

	verificationToken, err := GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Create PasswordService instance
	passwordService := PasswordService.NewPasswordService()

	// Call method on instance
	hashedPassword, err := passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UUID:                       u,
		Username:                   username,
		Fullname:                   fullname,
		Email:                      email,
		Password:                   hashedPassword,
    		ResetPasswordToken:          "",
		ResetPasswordTokenCreatedAt: nil,
		AccountVerificationToken:   &signupToken,
		VerificationTokens:         verificationToken,
		VerificationTokenCreatedAt: time.Now(),
		IsVerified:                 false,
		CreatedAt:                  time.Now(),
		TwoFASecret:                nil,
		TwoFAEnabled:               false,
	}
	return user, nil
}

// generateUUID return new UUID
func generateUUID() string {
	return uuid.New().String()
}