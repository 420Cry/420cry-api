package oauth

import (
	Email "cry-api/app/email"
	Repository "cry-api/app/repositories"
	EmailService "cry-api/app/services/email"
	OAuthService "cry-api/app/services/oauth"
	PasswordService "cry-api/app/services/password"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

type OAuthController struct {
	OAuthService OAuthService.OAuthServiceInterface
	UserService  UserService.UserServiceInterface
}

func NewOAuthController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *OAuthController {
	userRepository := Repository.NewGormUserRepository(db)
	oauthRepository := Repository.NewGormOAuthRepository(db)

	passwordService := PasswordService.NewPasswordService()
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)

	emailService := EmailService.NewEmailService(emailSender)
	authService := UserService.NewAuthService(userRepository, passwordService)
	verificationService := UserService.NewVerificationService(userRepository)

	OAuthService := OAuthService.NewOAuthService(oauthRepository)
	userService := UserService.NewUserService(userRepository, emailService, verificationService, authService)

	return &OAuthController{OAuthService: OAuthService, UserService: userService}
}
