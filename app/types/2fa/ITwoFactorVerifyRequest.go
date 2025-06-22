// Package types provides types
package types

// ITwoFactorVerifyRequest represents the request payload for 2FA setup/verification.
type ITwoFactorVerifyRequest struct {
	UserUUID string `json:"userUUID"` // UUID should be required in practice to identify user
	OTP      string `json:"otp"`      // OTP
}
