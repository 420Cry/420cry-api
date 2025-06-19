// Package services provides utility functions for password hashing and verification.
package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService provides password hashing and verification methods.
type PasswordService struct{}

// PasswordServiceInterface defines methods for password operations.
type PasswordServiceInterface interface {
	HashPassword(password string) (string, error)
	CheckPassword(hashed, plain string) error
}

// NewPasswordService returns a new PasswordService instance.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword takes a plain-text password and returns its bcrypt hash.
// It is used to securely store user passwords.
func (p *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash error: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a hashed password with a plain-text password.
// It returns nil if the passwords match, or an error if they do not.
func (p *PasswordService) CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
