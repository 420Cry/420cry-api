// Package services provides business logic for handling email operations.
package services

import (
	"log"

	Email "cry-api/app/email"
	"cry-api/app/utils"
)

// EmailServiceInterface provides all EmailService methods
type EmailServiceInterface interface {
	SendVerifyAccountEmail(to, from, username, verificationLink, verificationToken string) error
	SendResetPasswordEmail(to, from, username, resetPasswordLink, APIURL string) error
	SendTwoFactorAlternativeEmail(to, from, username, otp string, expiryMinutes int) error
}

// EmailSender is an interface for sending emails
type EmailSender interface {
	Send(email Email.EmailMessage) error
}

// EmailCreator is an interface for creating email templates
type EmailCreator interface {
	CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken string) (Email.EmailMessage, error)
	CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (Email.EmailMessage, error)
	CreateTwoFactorAlternativeEmail(to, from, userName, otp string, expiryMinutes int) (Email.EmailMessage, error)
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
func (service *EmailService) SendVerifyAccountEmail(to, from, userName, verificationLink, verificationToken string) error {
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)
	verificationLink = utils.SanitizeInput(verificationLink)

	email, err := service.emailCreator.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken)
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

	return nil
}

// SendTwoFactorAlternativeEmail creates the alternative 2FA OTP email and sends it
func (service *EmailService) SendTwoFactorAlternativeEmail(to, from, userName, otp string, expiryMinutes int) error {
	to = utils.SanitizeInput(to)
	userName = utils.SanitizeInput(userName)

	email, err := service.emailCreator.CreateTwoFactorAlternativeEmail(to, from, userName, otp, expiryMinutes)
	if err != nil {
		log.Printf("Error creating 2FA email template: %v", err)
		return err
	}

	err = service.emailSender.Send(email)
	if err != nil {
		log.Printf("Error sending 2FA email: %v", err)
		return err
	}

	log.Printf("2FA alternative email sent successfully to %s", email.To)
	return nil
}
