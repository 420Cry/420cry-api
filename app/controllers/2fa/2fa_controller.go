// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"cry-api/app/container"
	TwoFactorService "cry-api/app/services/2fa"
	AuthService "cry-api/app/services/auth"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
)

// TwoFactorController handles 2FA-related HTTP requests.
type TwoFactorController struct {
	UserService      UserService.UserServiceInterface
	UserTokenService UserService.UserTokenServiceInterface
	AuthService      AuthService.AuthServiceInterface
	TwoFactorService TwoFactorService.TwoFactorServiceInterface
	EmailService     EmailService.EmailServiceInterface
}

// NewTwoFactorController initializes a new TwoFactorController with dependencies from the container.
func NewTwoFactorController(container *container.Container) *TwoFactorController {
	// Get services from container
	userService := container.GetUserService()
	userTokenService := container.GetUserTokenService()
	authService := container.GetAuthService()
	emailService := container.GetEmailService()

	// Create 2FA service
	twoFactorService := TwoFactorService.NewTwoFactorService()

	return &TwoFactorController{
		UserService:      userService,
		UserTokenService: userTokenService,
		AuthService:      authService,
		TwoFactorService: twoFactorService,
		EmailService:     emailService,
	}
}
