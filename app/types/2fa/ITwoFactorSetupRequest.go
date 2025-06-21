// Package types provides types
package types

// ITwoFactorSetupRequest represents the request payload for 2FA setup/verification.
type ITwoFactorSetupRequest struct {
	UserUUID string  `json:"uuid"`             // UUID should be required in practice to identify user
	OTP      *string `json:"otp,omitempty"`    // OTP optional
	Secret   *string `json:"secret,omitempty"` // Secret optional
}
