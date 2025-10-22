// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"cry-api/app/container"
	AuthService "cry-api/app/services/auth"
	PasswordService "cry-api/app/services/auth/password"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	AuthService      AuthService.AuthServiceInterface
	UserService      UserService.UserServiceInterface
	EmailService     EmailService.EmailServiceInterface
	PasswordService  PasswordService.PasswordServiceInterface
	UserTokenService UserService.UserTokenServiceInterface
}

/*
NewUserController initializes and returns a new UserController instance with its dependencies from the container.
*/
func NewUserController(container *container.Container) *UserController {
	return &UserController{
		UserService:      container.GetUserService(),
		UserTokenService: container.GetUserTokenService(),
		EmailService:     container.GetEmailService(),
		AuthService:      container.GetAuthService(),
		PasswordService:  container.GetPasswordService(),
	}
}
