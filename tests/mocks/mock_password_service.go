package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockPasswordService mocks PasswordServiceInterface
type MockPasswordService struct {
	mock.Mock
}

// HashPassword mocks HashPassword from PasswordService
func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

// CheckPassword mocks CheckPassword from PasswordService
func (m *MockPasswordService) CheckPassword(hashed, plain string) error {
	args := m.Called(hashed, plain)
	return args.Error(0)
}
