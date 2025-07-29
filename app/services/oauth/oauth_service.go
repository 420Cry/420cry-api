package oauth

import (
	"context"
	"cry-api/app/config"
	"cry-api/app/factories"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	ConstantStore "cry-api/app/constants"
	Models "cry-api/app/models"
	Repositories "cry-api/app/repositories"
	OAuthType "cry-api/app/types/oauth"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type OAuthService struct {
	OAuthRepo Repositories.OAuthRepository
}

type OAuthServiceInterface interface {
	ExchangeToken(c context.Context, code string) (*oauth2.Token, error)
	FetchUserInfo(token *oauth2.Token) (*OAuthType.IGoogleUserResponse, error)
	CreateGoogleAccount(existingUser *Models.User, googleUserInfo *OAuthType.IGoogleUserResponse, token *oauth2.Token) error
	FindAccountByProviderAndId(provider, providerId string) (*Models.Oauth_Accounts, error)
}

func NewOAuthService(oauthRepo Repositories.OAuthRepository) *OAuthService {
	return &OAuthService{OAuthRepo: oauthRepo}
}

func (s *OAuthService) ExchangeToken(c context.Context, code string) (*oauth2.Token, error) {
	googleCfg := config.GetOAuthConfig("google")
	token, err := googleCfg.Exchange(c, code)

	if err != nil {
		log.Println("cannot exchange token")
		return nil, err
	}

	return token, nil
}

func (s *OAuthService) FetchUserInfo(token *oauth2.Token) (*OAuthType.IGoogleUserResponse, error) {
	googleUserApi, err := ConstantStore.ReturnConstant("googleUserApi")
	if err != nil {
		return nil, fmt.Errorf("cannot find google api")
	}

	req, err := http.NewRequest(http.MethodGet, googleUserApi, nil)
	if err != nil {
		return nil, fmt.Errorf("failed creating request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("cannot make request to google server")
	}

	defer resp.Body.Close()

	var googleUserInfo *OAuthType.IGoogleUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, fmt.Errorf("google response is invalid")
	}

	return googleUserInfo, nil
}

func (s *OAuthService) CreateGoogleAccount(existingUser *Models.User, googleUserInfo *OAuthType.IGoogleUserResponse, token *oauth2.Token) error {
	provider := "Google"
	providerId := googleUserInfo.Sub
	oauthAccount, err := factories.NewOAuthAccount(existingUser, provider, providerId, googleUserInfo.Email, token)

	if err != nil {
		log.Println("cannot create oauth account", err)
	}

	if err := s.OAuthRepo.Save(oauthAccount); err != nil {
		log.Println("cannot save oauth account:", err)
		return err
	}

	return nil
}

func (s *OAuthService) FindAccountByProviderAndId(provider, providerId string) (*Models.Oauth_Accounts, error) {
	oauthAccount, err := s.OAuthRepo.FindByProviderAndId(provider, providerId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("record not found")
		}

		return nil, fmt.Errorf("error finding oauth account: %s", err)
	}

	return oauthAccount, nil
}
