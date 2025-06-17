package mocks

import (
	UserModel "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockAuthService mocks AuthService methods
type MockAuthService struct {
	mock.Mock
}

// AuthenticateUser mocks AuthenticateUser method from AuthService
func (m *MockAuthService) AuthenticateUser(username, password string) (*UserModel.User, error) {
	args := m.Called(username, password)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}
