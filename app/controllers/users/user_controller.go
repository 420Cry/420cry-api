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
	AuthService      UserServices.AuthServiceInterface
	UserService      UserServices.UserServiceInterface
	EmailService     EmailServices.EmailServiceInterface
	PasswordService  PasswordService.PasswordServiceInterface
	UserTokenService UserServices.UserTokenServiceInterface
}

/*
NewUserController initializes and returns a new NewUserController instance with its dependencies.
*/
func NewUserController(db *gorm.DB, cfg *EnvTypes.EnvConfig) *UserController {
	passwordService := PasswordService.NewPasswordService()
	userRepository := UserRepository.NewGormUserRepository(db)
	userTokenRepository := UserRepository.NewGormUserTokenRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)

	emailCreator := &EmailServices.EmailCreatorImpl{}
	emailService := EmailServices.NewEmailService(emailSender, emailCreator)

	authService := UserServices.NewAuthService(userRepository, passwordService)
	userTokenService := UserServices.NewUserTokenService(userTokenRepository)
	userService := UserServices.NewUserService(
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
