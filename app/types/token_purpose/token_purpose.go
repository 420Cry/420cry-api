// Package types defines constants and types used across the application.
package types

// TokenPurpose represents the purpose of a user token.
type TokenPurpose string

const (
	// AccountVerification is used for long-link account verification tokens
	// sent when a user creates a new account. This token usually arrives via email.
	AccountVerification TokenPurpose = "account_verification"

	// AccountVerificationOTP is used for one-time passcodes (OTP) during
	// account verification, typically sent alongside or after the long-link token.
	AccountVerificationOTP TokenPurpose = "account_verification_otp"

	// ResetPassword is used for password reset tokens. These tokens allow
	// a verified user to reset their password securely.
	ResetPassword TokenPurpose = "reset_password"

	// TwoFactorAuthAlternativeOTP is used for two-factor authentication (2FA) alternative (when users can not access to their auth app),
	TwoFactorAuthAlternativeOTP TokenPurpose = "two_factor_auth_alternative_otp"
)
