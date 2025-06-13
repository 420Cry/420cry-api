// Package types provides type definitions for user signup requests.
package types

// TokenRequest represents the request payload containing the verification token sent to a user's email during account signup.
type UserVerifyAccountTokenRequest struct {
	Token string `json:"token"`
}
