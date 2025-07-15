// Package services provides business logic for handling email operations.
package services

import (
	"log"

	Email "cry-api/app/email"
	"cry-api/app/utils"
)

// EmailServiceInterface provides all EmailService methods
type EmailServiceInterface interface {
	SendVerifyAccountEmail(to, from, username, verificationLink, verificationTokens string) error
	SendResetPasswordEmail(to, from, username, resetPasswordLink, APIURL string) error
}

// EmailSender is an interface for sending emails
type EmailSender interface {
	Send(email Email.EmailMessage) error
}

// EmailCreator is an interface for creating email templates
type EmailCreator interface {
	CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens string) (Email.EmailMessage, error)
	CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (Email.EmailMessage, error)
}

// EmailService provides operations for sending emails
type EmailService struct {
	emailSender  EmailSender
	emailCreator EmailCreator
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(emailSender EmailSender, emailCreator EmailCreator) *EmailService {
	return &EmailService{
		emailSender:  emailSender,
		emailCreator: emailCreator,
	}
}

// SendVerifyAccountEmail creates an email and sends it
func (service *EmailService) SendVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens string) error {
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)
	verificationLink = utils.SanitizeInput(verificationLink)

	email, err := service.emailCreator.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens)
	if err != nil {
		return err
	}

	return service.emailSender.Send(email)
}

// SendResetPasswordEmail creates the reset password email and sends it
func (service *EmailService) SendResetPasswordEmail(to, from, userName, resetPasswordLink, APIURL string) error {
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)
	resetPasswordLink = utils.SanitizeInput(resetPasswordLink)

	email, err := service.emailCreator.CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL)
	if err != nil {
		log.Printf("Error creating email template: %v", err)
	}

	err = service.emailSender.Send(email)
	if err != nil {
		log.Printf("error sending the email: %v", err)
	}

	log.Printf("Email sent successfully to %s", email.To)
	return nil
}
