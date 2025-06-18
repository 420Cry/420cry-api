package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	EmailServices "cry-api/app/services/email"
	UserServices "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// TwoFactorController handles 2FA-related HTTP requests.
type TwoFactorController struct {
	UserService UserServices.UserServiceInterface
}

// NewTwoFactorController initializes a new TwoFactorController with dependencies.
func NewTwoFactorController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *TwoFactorController {
	userRepo := UserRepository.NewGormUserRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)
	userService := UserServices.NewUserService(userRepo, emailService)

	return &TwoFactorController{
		UserService: userService,
	}
}
