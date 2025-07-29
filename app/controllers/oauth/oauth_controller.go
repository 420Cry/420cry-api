package oauth

import (
	Email "cry-api/app/email"
	Repository "cry-api/app/repositories"
	EmailServices "cry-api/app/services/email"
	OAuthService "cry-api/app/services/oauth"
	PasswordService "cry-api/app/services/password"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// OAuthController handles OAuth HTTP Requests
type OAuthController struct {
	OAuthService OAuthService.OAuthServiceInterface
	UserService  UserService.UserServiceInterface
}

// NewOAuthController initializes OAuthController to be used for OAuth route handler
func NewOAuthController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *OAuthController {
	userRepository := Repository.NewGormUserRepository(db)
	oauthRepository := Repository.NewGormOAuthRepository(db)
	transactionRepository := Repository.NewGormTransactionRepository(db)

	passwordService := PasswordService.NewPasswordService()

	emailCreator := &EmailServices.EmailCreatorImpl{}
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)

	emailService := EmailServices.NewEmailService(emailSender, emailCreator)
	authService := UserService.NewAuthService(userRepository, passwordService)
	verificationService := UserService.NewVerificationService(userRepository)

	OAuthService := OAuthService.NewOAuthService(oauthRepository)
	userService := UserService.NewUserService(userRepository, transactionRepository, emailService, verificationService, authService)

	return &OAuthController{OAuthService: OAuthService, UserService: userService}
}
