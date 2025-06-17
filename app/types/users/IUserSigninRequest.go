// Package types provides type definitions for user signup requests.
package types

// IUserSigninRequest represents the payload required for a user to sign in
type IUserSigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
