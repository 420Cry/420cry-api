package emaildomain

// CreateVerifyAccountEmail generates an email template for verifying an account
func CreateVerifyAccountEmail(to, userName, verificationLink string) (Email, error) {
	// Prepare the template data
	data := map[string]interface{}{
		"UserName":         userName,
		"AppName":          "420Cry",
		"VerificationLink": verificationLink,
		"Year":             "2025",
	}

	// Render the HTML body (RenderHTMLTemplate should be part of the infrastructure)
	htmlBody, err := RenderHTMLTemplate("app/templates/email/verify_account.html", data)
	if err != nil {
		return Email{}, err
	}

	return NewEmail(to, "Verify Your Account", htmlBody), nil
}
