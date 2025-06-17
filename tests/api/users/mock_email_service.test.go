// Package tests provides tests for user routes.
package tests

import (
	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock of EmailService for testing
type MockEmailService struct {
	mock.Mock
}

// SendVerifyAccountEmail mocks sending a verification email to a user with the provided recipient, sender, username, verification link, and token.
func (m *MockEmailService) SendVerifyAccountEmail(to, from, username, link, token string) error {
	args := m.Called(to, from, username, link, token)
	return args.Error(0)
}

// SendResetPasswordEmail mocks sending a reset password email to a user with provided recipient, sender, username and reset password link
func (m *MockEmailService) SendResetPasswordEmail(to, from, userName, resetPasswordLink string) error {
	args := m.Called(to, from, userName, resetPasswordLink)
	return args.Error(0)
}
