// Package services provides business logic for handling email operations.
package services

import (
	Email "cry-api/app/email"
	"cry-api/app/utils"
	"log"
)

// EmailServiceInterface provides all EmailService methods
type EmailServiceInterface interface {
	SendVerifyAccountEmail(to, from, username, verificationLink, verificationTokens string) error
	SendResetPasswordEmail(to, from, username, resetPasswordLink, APIURL string) error
}

// EmailService provides operations for sending emails
type EmailService struct {
	emailSender EmailSender
}

// EmailSender is an interface for sending emails
type EmailSender interface {
	Send(email Email.EmailMessage) error
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

	email, err := Email.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens)
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

// SendResetPasswordEmail creates the reset password email and send to the user
func (service *EmailService) SendResetPasswordEmail(to, from, userName, resetPasswordLink, APIURL string) error {
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)
	resetPasswordLink = utils.SanitizeInput(resetPasswordLink)

	// Creating email template
	email, err := Email.CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL)

	if err != nil {
		log.Printf("Error creating email template: %v", err)
	}

	// Sending the email
	err = service.emailSender.Send(email)

	if err != nil {
		log.Printf("error sending the email: %v", err)
	}

	log.Printf("Email sent successfully to %s", email.To)
	return nil
}
