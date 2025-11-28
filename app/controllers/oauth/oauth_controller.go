package oauth

import (
	"cry-api/app/container"
	OAuthService "cry-api/app/services/oauth"
	UserService "cry-api/app/services/users"
)

// OAuthController handles OAuth HTTP Requests
type OAuthController struct {
	OAuthService OAuthService.OAuthServiceInterface
	UserService  UserService.UserServiceInterface
}

// NewOAuthController initializes OAuthController to be used for OAuth route handler
func NewOAuthController(container *container.Container) *OAuthController {
	return &OAuthController{OAuthService: container.GetOAuthService(), UserService: container.GetUserService()}
}
