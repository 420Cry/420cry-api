package services

import (
	EmailDomain "cry-api/app/domain/email"
	"cry-api/app/utils"
)

// EmailService provides operations for sending emails
type EmailService struct {
	emailSender EmailSender
}

// EmailSender is an interface for sending emails
type EmailSender interface {
	Send(email EmailDomain.EmailMessage) error
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(emailSender EmailSender) *EmailService {
	return &EmailService{emailSender: emailSender}
}

// SendVerifyAccountEmail creates an email and sends it
func (service *EmailService) SendVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens string) error {
	// Sanitize inputs
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)
	verificationLink = utils.SanitizeInput(verificationLink)

	email, err := EmailDomain.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens)
	if err != nil {
		return err
	}

	// Send the email via the core layer
	err = service.emailSender.Send(email)
	if err != nil {
		return err
	}
	return nil
}
