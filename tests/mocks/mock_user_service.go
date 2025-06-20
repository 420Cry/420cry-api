package mocks

import (
	"github.com/stretchr/testify/mock"

	UserModel "cry-api/app/models"
)

// MockUserService mocks UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// CreateUser mocks CreateUser from UserService
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

// UpdateUser mocks UpdateUser from UserService
func (m *MockUserService) UpdateUser(user *UserModel.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindUserByEmail mocks FindUserByEmail from UserService
func (m *MockUserService) FindUserByEmail(email string) (*UserModel.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}

// FindUserByResetPasswordToken mocks FindUserByResetPasswordToken from userService
func (m *MockUserService) FindUserByResetPasswordToken(token string) (*UserModel.User, error) {
	args := m.Called(token)
	user, _ := args.Get(0).(*UserModel.User)
	return user, args.Error(1)
}
