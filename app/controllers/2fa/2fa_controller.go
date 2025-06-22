package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	TwoFactorService "cry-api/app/services/2fa"
	EmailServices "cry-api/app/services/email"
	PasswordService "cry-api/app/services/password"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// TwoFactorController handles 2FA-related HTTP requests.
type TwoFactorController struct {
	UserService      UserService.UserServiceInterface
	AuthService      UserService.AuthServiceInterface
	TwoFactorService TwoFactorService.TwoFactorServiceInterface
}

// NewTwoFactorController initializes a new TwoFactorController with dependencies.
func NewTwoFactorController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *TwoFactorController {
	passwordService := PasswordService.NewPasswordService()
	userRepository := UserRepository.NewGormUserRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailService := EmailServices.NewEmailService(emailSender)

	authService := UserService.NewAuthService(userRepository, passwordService)
	verificationService := UserService.NewVerificationService(userRepository)

	userService := UserService.NewUserService(userRepository, emailService, verificationService, authService)

	return &TwoFactorController{
		UserService: userService,
		AuthService: authService,
	}
}
