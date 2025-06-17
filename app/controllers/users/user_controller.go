// Package controllers contains HTTP handlers for user-related operations.
package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	EmailServices "cry-api/app/services/email"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	UserService  UserServices.UserServiceInterface
	EmailService EmailServices.EmailServiceInterface
}

/*
NewUserController initializes and returns a new NewUserController instance with its dependencies.
It sets up the user repository, email sender, email service, and user service
using the provided GORM database connection and environment configuration.

Parameters:
  - db:   A pointer to the GORM database instance.
  - cfg:  A pointer to the environment configuration.

Returns:
  - A pointer to the initialized Handler.
*/
func NewUserController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *UserController {
	userRepository := UserRepository.NewGormUserRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)
	userService := UserServices.NewUserService(userRepository, emailService)

	return &UserController{UserService: userService, EmailService: emailService}
}
