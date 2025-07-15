// Package services provides business logic for handling email operations.
package services

import Email "cry-api/app/email"

// EmailCreatorImpl implements the EmailCreator interface
type EmailCreatorImpl struct{}

// CreateVerifyAccountEmail creates the verification email
func (e *EmailCreatorImpl) CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens string) (Email.EmailMessage, error) {
	return Email.CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationTokens)
}

// CreateResetPasswordRequestEmail creates the reset password request email
func (e *EmailCreatorImpl) CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (Email.EmailMessage, error) {
	return Email.CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL)
}
