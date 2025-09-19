// Package types provides types
package types

// ITwoFactorAlternativeRequest represents the request payload for 2FA alternative verification.
type ITwoFactorAlternativeRequest struct {
	Email string `json:"email"`
}
