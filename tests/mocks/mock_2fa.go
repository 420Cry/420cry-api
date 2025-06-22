package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockTwoFactorService mocks the functions from the 2fa services package
type MockTwoFactorService struct {
	mock.Mock
}

// VerifyTOTP mocks VerifyTOTP for TwoFactorService
func (m *MockTwoFactorService) VerifyTOTP(secret, token string) bool {
	args := m.Called(secret, token)
	return args.Bool(0)
}

// GenerateTOTP mocks GenerateTOTP for TwoFactorService
func (m *MockTwoFactorService) GenerateTOTP(userEmail string) (string, string, error) {
	args := m.Called(userEmail)
	return args.String(0), args.String(1), args.Error(2)
}

// GenerateOtpauthURL mocks GenerateOtpauthURL for TwoFactorService
func (m *MockTwoFactorService) GenerateOtpauthURL(userEmail, secret string) string {
	args := m.Called(userEmail, secret)
	return args.String(0)
}

// GenerateQRCodeBase64 mocks GenerateQRCodeBase64 for TwoFactorService
func (m *MockTwoFactorService) GenerateQRCodeBase64(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}
