// Package mail provides functionality for creating and sending email messages,
// including templated emails for user account verification and other notifications.
package mail

import (
	"fmt"
	"time"

	"cry-api/app/utils"
)

// CreateVerifyAccountEmail generates an EmailMessage for verifying a user account.
// It populates the email template with the username, verification link, token, and current year.
//
// Parameters:
//   - to: recipient email address
//   - from: sender email address
//   - userName: recipient's username to personalize the email
//   - verificationLink: URL for account verification
//   - token: one-time verification token (OTP) to include in the email
//
// Returns:
//   - an EmailMessage with subject "Verify Your Account" and the rendered HTML body
//   - an error if the template rendering fails
func CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken string) (EmailMessage, error) {
	data := map[string]any{
		"UserName":         userName,
		"AppName":          "420Cry",
		"VerificationLink": verificationLink,
		"Token":            verificationToken,
		"Year":             time.Now().Year(),
	}

	templatePrefix := utils.GenerateEmailTemplatePrefix()

	templatePath := fmt.Sprintf("%s/verify_account.html", templatePrefix)

	htmlBody, err := RenderTemplate(templatePath, data)
	if err != nil {
		return EmailMessage{}, fmt.Errorf("template render error: %w", err)
	}

	return NewEmailMessage(to, from, "Verify Your Account", htmlBody), nil
}
