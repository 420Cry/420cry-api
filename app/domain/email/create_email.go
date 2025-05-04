package domain

import (
	"time"
)

func CreateVerifyAccountEmail(to, from, userName, verificationLink string) (EmailMessage, error) {
	data := map[string]any{
		"UserName":         userName,
		"AppName":          "420Cry",
		"VerificationLink": verificationLink,
		"Year":             time.Now().Year(),
	}
	htmlBody, err := RenderHTMLTemplate("app/app/templates/email/verify_account.html", data)

	if err != nil {
		return EmailMessage{}, err
	}

	return NewEmailMessage(to, from, "Verify Your Account", htmlBody), nil
}
