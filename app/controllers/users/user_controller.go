package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	EmailServices "cry-api/app/services/email"
	PasswordService "cry-api/app/services/password"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// UserController handles HTTP requests related to user operations
type UserController struct {
	VerificationService UserServices.VerificationServiceInterface
	AuthService         UserServices.AuthServiceInterface
	UserService         UserServices.UserServiceInterface
	EmailService        EmailServices.EmailServiceInterface
	PasswordService     PasswordService.PasswordServiceInterface
}

/*
NewUserController initializes and returns a new NewUserController instance with its dependencies.
*/
func NewUserController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *UserController {
	userRepository := UserRepository.NewGormUserRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)

	userService := UserServices.NewUserService(userRepository, emailService)
	return &UserController{
		UserService:  userService,
		EmailService: emailService,
	}
}
