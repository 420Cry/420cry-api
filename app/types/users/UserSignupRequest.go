// Package types provides type definitions for user signup requests.
package types

// UserSignupRequest represents the payload required for a user to sign up, including fullname, username, email, and password.
type UserSignupRequest struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
