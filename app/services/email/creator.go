// Package services provides business logic for handling email operations.
package services

import Email "cry-api/app/email"

// EmailCreatorImpl implements the EmailCreator interface
type EmailCreatorImpl struct{}

// CreateVerifyAccountEmail creates the verification email
func (e *EmailCreatorImpl) CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken string) (Email.EmailMessage, error) {
	return Email.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken)
}

// CreateResetPasswordRequestEmail creates the reset password request email
func (e *EmailCreatorImpl) CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (Email.EmailMessage, error) {
	return Email.CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL)
}

// CreateTwoFactorAlternativeEmail creates an alternative 2FA OTP email
func (e *EmailCreatorImpl) CreateTwoFactorAlternativeEmail(to, from, userName, otp string, expiryMinutes int) (Email.EmailMessage, error) {
	return Email.CreateTwoFactorAlternativeEmail(to, from, userName, otp, expiryMinutes)
}
