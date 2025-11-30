package mocks

import (
	Models "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

type MockOAuthTransactionRepository struct {
	mock.Mock
}

func (m *MockOAuthTransactionRepository) CreateUserByOAuth(user *Models.User, oauthAccount *Models.Oauth_Accounts) error {
	args := m.Called(user, oauthAccount)
	return args.Error(0)
}
