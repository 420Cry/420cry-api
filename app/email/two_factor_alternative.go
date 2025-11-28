// Package mail provides functionality for creating and sending email messages,
// including templated emails for 2FA alternative OTP delivery.
package mail

import (
	"fmt"
	"time"

	"cry-api/app/utils"
)

// CreateTwoFactorAlternativeEmail generates an EmailMessage for sending a one-time
// alternative 2FA code via email. No link is included — just the OTP.
//
// Parameters:
//   - to: recipient email address
//   - from: sender email address
//   - userName: recipient's username to personalize the email
//   - otp: one-time code to include in the email
//   - expiryMinutes: validity duration for the OTP (in minutes)
//
// Returns:
//   - an EmailMessage with subject "Your One-Time Verification Code"
//   - an error if the template rendering fails
func CreateTwoFactorAlternativeEmail(to, from, userName, otp string, expiryMinutes int) (EmailMessage, error) {
	data := map[string]any{
		"UserName":      userName,
		"AppName":       "420Cry",
		"Otp":           otp,
		"ExpiryMinutes": expiryMinutes,
		"Year":          time.Now().Year(),
	}

	templatePrefix := utils.GenerateEmailTemplatePrefix()
	templatePath := fmt.Sprintf("%s/two_factor_alternative.html", templatePrefix)

	htmlBody, err := RenderTemplate(templatePath, data)
	if err != nil {
		return EmailMessage{}, fmt.Errorf("template render error: %w", err)
	}

	return NewEmailMessage(to, from, "Your One-Time Verification Code", htmlBody), nil
}
