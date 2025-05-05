package domain

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type User struct {
	ID                 int       `json:"id"`
	UUID               string    `json:"uuid" gorm:"unique;not null"`
	Username           string    `json:"username" gorm:"unique;not null"`
	Email              string    `json:"email" gorm:"unique;not null"`
	Fullname           string    `json:"fullname"`
	Password           string    `json:"-" gorm:"not null"`
	Token              string    `json:"token,omitempty" gorm:"unique"`
	VerificationTokens string    `json:"verification_tokens,omitempty" gorm:"size:6"`
	IsVerified         bool      `json:"is_verified" gorm:"not null;default:false"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:timestamp;default:NULL;autoUpdateTime"`
}

// NewUser creates a new User entity with hashed password
func NewUser(fullname, username, email, password string) (*User, error) {
	uuid := uuid.New().String()
	signupToken, err := GenerateSignupToken()
	if err != nil {
		return nil, err
	}

	verificationToken, err := GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	// Hash the password before storing it
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		UUID:               uuid,
		Username:           username,
		Fullname:           fullname,
		Email:              email,
		Password:           hashedPassword,
		Token:              signupToken,
		VerificationTokens: verificationToken,
		IsVerified:         false,
		CreatedAt:          time.Now(),
	}

	return user, nil
}

// HashPassword hashes the user's password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
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

// GenerateVerificationToken generates a 6-character random verification token
func GenerateVerificationToken() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const tokenLength = 6
	token := make([]byte, tokenLength)

	randomBytes := make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	for i, b := range randomBytes {
		token[i] = charset[int(b)%len(charset)]
	}

	return string(token), nil
}
