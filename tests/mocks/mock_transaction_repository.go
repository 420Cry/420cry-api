package mocks

import (
	Models "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateUserByOAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error {
	args := m.Called(user, oauthAccount)
	return args.Error(0)
}
