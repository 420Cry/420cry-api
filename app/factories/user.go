// Package factories provides factory functions for creating and initializing
package factories

import (
	"fmt"
	"time"

	"cry-api/app/models"
	PasswordService "cry-api/app/services/password"

	"github.com/google/uuid"
)

// NewUser creates a new instance of models.User with hashed password
func NewUser(fullname, username, email, password string) (*models.User, error) {
	passwordService := PasswordService.NewPasswordService()

	hashedPassword, err := passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UUID:       uuid.New().String(),
		Username:   username,
		Fullname:   fullname,
		Email:      email,
		Password:   hashedPassword,
		IsVerified: false,
		CreatedAt:  time.Now(),
	}

	return user, nil
}

// TokenType determines the kind of token to generate
type TokenType int

const (
	// LongLink represents a long-lived token, typically used for account verification links.
	LongLink TokenType = iota

	// OTP represents a short-lived one-time password (OTP) token.
	OTP
)

// NewUserToken creates a new UserToken for a user based on the TokenType and expiration duration
func NewUserToken(userID int, purpose string, duration time.Duration, tokenType TokenType) (*models.UserToken, error) {
	var tokenValue string
	var err error

	switch tokenType {
	case LongLink:
		tokenValue, err = Generate32ByteToken()
	case OTP:
		tokenValue, err = GenerateOTP()
	default:
		return nil, fmt.Errorf("invalid token type")
	}

	if err != nil {
		return nil, err
	}

	return &models.UserToken{
		UserID:    userID,
		Token:     tokenValue,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(duration),
		Consumed:  false,
	}, nil
}
