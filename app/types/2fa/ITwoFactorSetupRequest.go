// Package types provides types
package types

// ITwoFactorSetupRequest represents the request payload for 2FA setup/verification.
// Fields are optional to support different request phases.
type ITwoFactorSetupRequest struct {
	UserUUID string  `json:"uuid"`          // UUID should be required in practice to identify user
	OTP      *string `json:"otp,omitempty"` // OTP optional
}
