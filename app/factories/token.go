// Package factories provides factory functions for creating and initializing
package factories

import (
	"crypto/rand"
	"fmt"
)

// Generate32ByteToken creates a secure random 32-byte token returned as a hex string.
// This token is typically embedded in account activation URLs.
func Generate32ByteToken() (string, error) {
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

	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	for i := range b {
		b[i] = charset[int(randBytes[i])%len(charset)]
	}
	return string(b), nil
}
