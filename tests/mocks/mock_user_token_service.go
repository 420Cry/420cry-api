package mocks

import (
	"cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserTokenService is a mock implementation of UserTokenServiceInterface
type MockUserTokenService struct {
	mock.Mock
}

// Save mocks Save from user_token_service
func (m *MockUserTokenService) Save(token *models.UserToken) error {
	args := m.Called(token)
	return args.Error(0)
}

// FindValidToken mocks FindValidToken from user_token_service
func (m *MockUserTokenService) FindValidToken(token, purpose string) (*models.UserToken, error) {
	args := m.Called(token, purpose)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserToken), args.Error(1)
	}
	return nil, args.Error(1)
}

// FindLatestValidToken mocks FindLatestValidToken from user_token_service
func (m *MockUserTokenService) FindLatestValidToken(userID int, purpose string) (*models.UserToken, error) {
	args := m.Called(userID, purpose)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserToken), args.Error(1)
	}
	return nil, args.Error(1)
}

// ConsumeToken mocks ConsumeToken from user_token_service
func (m *MockUserTokenService) ConsumeToken(userID int, token, purpose string) error {
	args := m.Called(userID, token, purpose)
	return args.Error(0)
}

// DeleteExpired mocks DeleteExpired from user_token_service
func (m *MockUserTokenService) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}
