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
