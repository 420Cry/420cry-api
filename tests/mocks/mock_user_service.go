package mocks

import (
	UserModel "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService mocks UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// CreateUser mocks CreateUser from  UserService
func (m *MockUserService) CreateUser(fullname, username, email, password string) (*UserModel.User, error) {
	args := m.Called(fullname, username, email, password)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// GetUserByUUID mocks GetUserByUUID from UserService
func (m *MockUserService) GetUserByUUID(uuid string) (*UserModel.User, error) {
	args := m.Called(uuid)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}
