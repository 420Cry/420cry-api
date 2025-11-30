package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

func GetOAuthConfig(provider string) *oauth2.Config {
	if provider == "google" {
		return &oauth2.Config{
			ClientID:     Get().GoogleClientId,
			ClientSecret: Get().GoogleClientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{people.UserinfoEmailScope, people.UserinfoProfileScope},
			RedirectURL:  Get().GoogleRedirectUrl,
		}
	}

	// Discord oauth config
	return &oauth2.Config{}
}
