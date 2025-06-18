// Package types provides type definitions for user signup requests.
package types

// ITwoFactorSetupRequest represents the request payload containing the user uuid
type ITwoFactorSetupRequest struct {
	UserUUID string `json:"uuid"`
}
