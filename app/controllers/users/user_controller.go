package controllers

import (
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	AuthService "cry-api/app/services/auth"
	PasswordService "cry-api/app/services/auth"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
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
NewUserController initializes and returns a new NewUserController instance with its dependencies.
*/
func NewUserController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *UserController {
	passwordService := PasswordService.NewPasswordService()
	userRepository := UserRepository.NewGormUserRepository(db)
	userTokenRepository := UserRepository.NewGormUserTokenRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)

	emailCreator := &EmailService.EmailCreatorImpl{}
	emailService := EmailService.NewEmailService(emailSender, emailCreator)

	authService := AuthService.NewAuthService(userRepository, passwordService)
	userTokenService := UserService.NewUserTokenService(userTokenRepository)
	userService := UserService.NewUserService(
		userRepository,
		userTokenRepository,
		emailService,
		authService,
	)

	return &UserController{
		UserService:      userService,
		UserTokenService: userTokenService,
		EmailService:     emailService,
		AuthService:      authService,
		PasswordService:  passwordService,
	}
}
