// Package mocks provides mocks for testing
package mocks

import (
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
