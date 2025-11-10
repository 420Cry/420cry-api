// Package types provides type definitions for user settings update requests.
package types

// IUserUpdateAccountNameRequest represents the payload required for updating a user's account name (fullname).
type IUserUpdateAccountNameRequest struct {
	AccountName string `json:"username" binding:"required"`
}
