package mocks

import (
	Models "cry-api/app/models"

	"github.com/stretchr/testify/mock"
)

type MockOAuthRepository struct {
	mock.Mock
}

func (m *MockOAuthRepository) Save(oauth_account *Models.Oauth_Accounts) error {
	args := m.Called(oauth_account)
	return args.Error(0)
}

func (m *MockOAuthRepository) FindByProviderAndId(provider, providerId string) (*Models.Oauth_Accounts, error) {
	args := m.Called(provider, providerId)
	oauthAccount := args.Get(0)

	if oauthAccount == nil {
		return nil, args.Error(1)
	}

	return oauthAccount.(*Models.Oauth_Accounts), nil
}
