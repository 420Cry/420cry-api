package mocks

import (
	"cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserTokenRepository mocks UserTokenRepository interface.
type MockUserTokenRepository struct {
	mock.Mock
}

// Save mocks Save method
func (m *MockUserTokenRepository) Save(token *models.UserToken) error {
	args := m.Called(token)
	return args.Error(0)
}

// FindValidToken mocks FindValidToken method
func (m *MockUserTokenRepository) FindValidToken(token, purpose string) (*models.UserToken, error) {
	args := m.Called(token, purpose)
	t := args.Get(0)
	if t == nil {
		return nil, args.Error(1)
	}
	return t.(*models.UserToken), args.Error(1)
}

// ConsumeToken mocks ConsumeToken method
func (m *MockUserTokenRepository) ConsumeToken(userID int, token, purpose string) error {
	args := m.Called(userID, token, purpose)
	return args.Error(0)
}

// DeleteExpired mocks DeleteExpired method
func (m *MockUserTokenRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// FindLatestValidToken mocks FindLatestValidToken method
func (m *MockUserTokenRepository) FindLatestValidToken(userID int, purpose string) (*models.UserToken, error) {
	args := m.Called(userID, purpose)
	t := args.Get(0)
	if t == nil {
		return nil, args.Error(1)
	}
	return t.(*models.UserToken), args.Error(1)
}
