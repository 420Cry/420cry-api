// Package domain contains business logic and domain models for the application.
package domain

import (
	"time"
)

// CreateVerifyAccountEmail generates an email message for account verification.
func CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens string) (EmailMessage, error) {
	data := map[string]any{
		"UserName":           userName,
		"AppName":            "420Cry",
		"VerificationLink":   verificationLink,
		"VerificationTokens": []string{verificationTokens},
		"Year":               time.Now().Year(),
	}

	// Render the HTML body with the template and data
	htmlBody, err := RenderHTMLTemplate("app/app/templates/email/verify_account.html", data)
	if err != nil {
		return EmailMessage{}, err
	}

	return NewEmailMessage(to, from, "Verify Your Account", htmlBody), nil
}

// CreateResetPasswordRequestEmail generates the email message for reset password 
func CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink string) (EmailMessage, error) {
	data := map[string]any{
		"UserName":          userName,
		"AppName":           "420Cry",
		"ResetPasswordLink": resetPasswordLink,
		"Year":              time.Now().Year(),
	}

	htmlBody, err := RenderHTMLTemplate("app/app/templates/email/reset_password.html", data)
	if err != nil {
		return EmailMessage{}, err
	}

	return NewEmailMessage(to, from, "Reset Your Password", htmlBody), nil

}
