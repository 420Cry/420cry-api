package factories

import (
	"cry-api/app/config"
	Models "cry-api/app/models"
	EncryptService "cry-api/app/services/encrypt"
	"log"

	"golang.org/x/oauth2"
)

func NewOAuthAccount(existingUser *Models.User, provider string, providerId string, token *oauth2.Token) (*Models.Oauth_Accounts, error) {
	oauthEncryptedKey := config.Get().OAuthEncryptedKey
	encryptedService := EncryptService.NewEncryptService()

	hashedAccessToken, err := encryptedService.EncryptToken([]byte(token.AccessToken), oauthEncryptedKey)
	if err != nil {
		log.Println("Cannot hash access token", err)
		return nil, err
	}

	hashedRefreshToken, err := encryptedService.EncryptToken([]byte(token.RefreshToken), oauthEncryptedKey)

	if err != nil {
		log.Println("Cannot hash refresh token", err)
		return nil, err
	}

	oauth_account := &Models.Oauth_Accounts{
		UserId:       existingUser.ID,
		Provider:     provider,
		ProviderId:   providerId,
		AccessToken:  hashedAccessToken,
		RefreshToken: hashedRefreshToken,
		TokenExpiry:  token.Expiry,
	}
	return oauth_account, nil
}
