package core

import "context"

type TwoFAService interface {
	// GenerateSecret creates a new 2FA secret and returns otpauth URL
	GenerateSecret(ctx context.Context, userID int, username, issuer string) (secret, otpURL string, err error)

	// VerifyCode verifies a TOTP code for a user
	VerifyCode(ctx context.Context, userID int, code string) (bool, error)

	// Enable2FA enables 2FA for user with given secret
	Enable2FA(ctx context.Context, userID int, secret string) error

	// Disable2FA disables 2FA for user
	Disable2FA(ctx context.Context, userID int) error
}
