package domain

import (
	"github.com/pquerna/otp/totp"
)

type TwoFA struct {
	Secret  string
	Enabled bool
}

// NewTwoFA creates a TwoFA instance from secret and enabled flag
func NewTwoFA(secret string, enabled bool) *TwoFA {
	return &TwoFA{
		Secret:  secret,
		Enabled: enabled,
	}
}

// VerifyCode checks the given TOTP code against the secret with some allowed skew
func (t *TwoFA) VerifyCode(code string) bool {
	return totp.Validate(code, t.Secret)
}

// Enable enables 2FA and sets the secret
func (t *TwoFA) Enable(secret string) {
	t.Secret = secret
	t.Enabled = true
}

// Disable disables 2FA (clear secret and flag)
func (t *TwoFA) Disable() {
	t.Secret = ""
	t.Enabled = false
}

// IsEnabled returns true if 2FA is enabled
func (t *TwoFA) IsEnabled() bool {
	return t.Enabled && t.Secret != ""
}
