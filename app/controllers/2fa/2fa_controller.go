package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	TwoFactorService "cry-api/app/services/2fa"
	AuthService "cry-api/app/services/auth"
	PasswordService "cry-api/app/services/auth"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// TwoFactorController handles 2FA-related HTTP requests.
type TwoFactorController struct {
	UserService      UserService.UserServiceInterface
	UserTokenService UserService.UserTokenServiceInterface
	AuthService      AuthService.AuthServiceInterface
	TwoFactorService TwoFactorService.TwoFactorServiceInterface
	EmailService     EmailService.EmailServiceInterface
}

// NewTwoFactorController initializes a new TwoFactorController with dependencies.
func NewTwoFactorController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *TwoFactorController {
	passwordService := PasswordService.NewPasswordService()
	userRepository := UserRepository.NewGormUserRepository(db)
	userTokenRepository := UserRepository.NewGormUserTokenRepository(db)

	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailCreator := &EmailService.EmailCreatorImpl{}
	emailService := EmailService.NewEmailService(emailSender, emailCreator)

	authService := AuthService.NewAuthService(userRepository, passwordService)
	userService := UserService.NewUserService(userRepository, userTokenRepository, emailService, authService)

	twoFactorService := TwoFactorService.NewTwoFactorService()
	userTokenService := UserService.NewUserTokenService(userTokenRepository)

	return &TwoFactorController{
		UserService:      userService,
		UserTokenService: userTokenService,
		AuthService:      authService,
		TwoFactorService: twoFactorService,
		EmailService:     emailService,
	}
}
