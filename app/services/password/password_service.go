// Package services provides utility functions for password hashing and verification.
package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain-text password and returns its bcrypt hash.
// It is used to securely store user passwords.
func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns the bcrypt hash of the password at the default cost.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash error: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a hashed password with a plain-text password.
// It returns nil if the passwords match, or an error if they do not.
func CheckPassword(hashed, plain string) error {
	// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
