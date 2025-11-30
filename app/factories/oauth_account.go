package factories

import (
	"cry-api/app/config"
	Models "cry-api/app/models"
	EncryptService "cry-api/app/services/encrypt"
	"log"

	"golang.org/x/oauth2"
)

func NewOAuthAccount(existingUser *Models.User, provider string, providerId string, email string, token *oauth2.Token) (*Models.Oauth_Accounts, error) {
	oauthEncryptedKey := config.Get().OAuthEncryptedKey
	encryptedService := EncryptService.NewEncryptService()

	encryptedAccessToken, err := encryptedService.EncryptToken([]byte(token.AccessToken), oauthEncryptedKey)
	if err != nil {
		log.Println("Cannot hash access token", err)
		return nil, err
	}

	encryptedRefreshToken, err := encryptedService.EncryptToken([]byte(token.RefreshToken), oauthEncryptedKey)

	if err != nil {
		log.Println("Cannot hash refresh token", err)
		return nil, err
	}

	oauth_account := &Models.Oauth_Accounts{
		UserId:       existingUser.ID,
		Email:        email,
		Provider:     provider,
		ProviderId:   providerId,
		AccessToken:  encryptedAccessToken,
		RefreshToken: encryptedRefreshToken,
		TokenExpiry:  token.Expiry,
	}
	return oauth_account, nil
}
