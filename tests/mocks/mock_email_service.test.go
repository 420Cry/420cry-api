// Package mocks provides mocks for testing
package mocks

import (
	Email "cry-api/app/email"

	"github.com/stretchr/testify/mock"
)

// MockEmailService mocks EmailServiceInterface
type MockEmailService struct {
	mock.Mock
}

// SendVerifyAccountEmail mocks SendVerifyAccountEmail from EmailService
func (m *MockEmailService) SendVerifyAccountEmail(to, from, username, link, token string) error {
	args := m.Called(to, from, username, link, token)
	return args.Error(0)
}

// SendResetPasswordEmail mocks SendResetPasswordEmail from EmailService
func (m *MockEmailService) SendResetPasswordEmail(to, from, username, resetPasswordLink, APIURL string) error {
	args := m.Called(to, from, username, resetPasswordLink, APIURL)
	return args.Error(0)
}

// SendTwoFactorAlternativeEmail mocks SendTwoFactorAlternativeEmail from EmailService
func (m *MockEmailService) SendTwoFactorAlternativeEmail(to, from, username, otp string, expiryMinutes int) error {
	args := m.Called(to, from, username, otp, expiryMinutes)
	return args.Error(0)
}

// MockEmailSender mocks the EmailSender interface
type MockEmailSender struct {
	mock.Mock
}

// Send mocks Send method of EmailSender
func (m *MockEmailSender) Send(email Email.EmailMessage) error {
	args := m.Called(email)
	return args.Error(0)
}

// MockEmailCreator mocks the EmailCreator interface
type MockEmailCreator struct {
	mock.Mock
}

// CreateVerifyAccountEmail mocks CreateVerifyAccountEmail from EmailCreator
func (m *MockEmailCreator) CreateVerifyAccountEmail(to, from, userName, verificationLink, verificationToken string) (Email.EmailMessage, error) {
	args := m.Called(to, from, userName, verificationLink, verificationToken)
	return args.Get(0).(Email.EmailMessage), args.Error(1)
}

// CreateResetPasswordRequestEmail mocks CreateResetPasswordRequestEmail from EmailCreator
func (m *MockEmailCreator) CreateResetPasswordRequestEmail(to, from, userName, resetPasswordLink, APIURL string) (Email.EmailMessage, error) {
	args := m.Called(to, from, userName, resetPasswordLink, APIURL)
	return args.Get(0).(Email.EmailMessage), args.Error(1)
}

// CreateTwoFactorAlternativeEmail mocks CreateTwoFactorAlternativeEmail from EmailCreator
func (m *MockEmailCreator) CreateTwoFactorAlternativeEmail(to, from, userName, otp string, expiryMinutes int) (Email.EmailMessage, error) {
	args := m.Called(to, from, userName, otp, expiryMinutes)
	return args.Get(0).(Email.EmailMessage), args.Error(1)
}
