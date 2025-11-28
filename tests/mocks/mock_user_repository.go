package mocks

import (
	"cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository mocks UserRepository interface.
type MockUserRepository struct {
	mock.Mock
}

// Save mocks Save method from UserRepository
func (m *MockUserRepository) Save(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// FindByUUID mocks FindByUUID method from UserRepository
func (m *MockUserRepository) FindByUUID(uuid string) (*models.User, error) {
	args := m.Called(uuid)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByID mocks FindByID method from UserRepository
func (m *MockUserRepository) FindByID(id int) (*models.User, error) {
	args := m.Called(id)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByEmail mocks FindByEmail method from UserRepository
func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByUsernameOrEmail mocks FindByUsernameOrEmail method from UserRepository
func (m *MockUserRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	args := m.Called(username, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByUsername mocks FindByUsername method from UserRepository
func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByAccountVerificationToken mocks FindByAccountVerificationToken method
func (m *MockUserRepository) FindByAccountVerificationToken(token string) (*models.User, error) {
	args := m.Called(token)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// FindByResetPasswordToken mocks FindByResetPasswordToken method
func (m *MockUserRepository) FindByResetPasswordToken(token string) (*models.User, error) {
	args := m.Called(token)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*models.User), args.Error(1)
}

// Delete mocks Delete method from UserRepository
func (m *MockUserRepository) Delete(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}
