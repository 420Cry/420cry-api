package mail

import (
	"cry-api/app/utils"
	"fmt"
	"time"
)

// CreateResetPasswordRequestEmail generates an EmailMessage for reset the user's password.
// It populates the email template with the username, reset password link, APIURL for logo and current year.
//
// Parameters:
//   - to: recipient email address
//   - from: sender email address
//   - userName: recipient's username to personalize the email
//   - resetPasswordLink: URL for reset password
//   - APIURL: Serve app logo
//
// Returns:
//   - an EmailMessage with subject "Reset your password" and the rendered HTML body
//   - an error if the template rendering fails
func CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (EmailMessage, error) {
	data := map[string]any{
		"UserName":          userName,
		"AppName":           "420Cry",
		"ResetPasswordLink": resetPasswordLink,
		"APIURL":            APIURL,
		"Year":              time.Now().Year(),
	}

	// Render a template path based on the environment
	templatePrefix := utils.GenerateEmailTemplatePrefix()

	templatePath := fmt.Sprintf("%s/reset_password.html", templatePrefix)
	htmlBody, err := RenderTemplate(templatePath, data)
	if err != nil {
		return EmailMessage{}, fmt.Errorf("template render error: %w", err)
	}

	return NewEmailMessage(to, from, "Reset Your Password", htmlBody), nil
}
