package controllers

import (
	Email "cry-api/app/email"
	Repository "cry-api/app/repositories"
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
	passwordService := PasswordService.NewPasswordService()
	userRepository := Repository.NewGormUserRepository(db)
	transactionRepository := Repository.NewGormTransactionRepository(db)
	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)

	// Instantiate EmailCreator implementation
	emailCreator := &EmailServices.EmailCreatorImpl{}

	// Pass both sender and creator
	emailService := EmailServices.NewEmailService(emailSender, emailCreator)

	authService := UserServices.NewAuthService(userRepository, passwordService)
	verificationService := UserServices.NewVerificationService(userRepository)
	userService := UserServices.NewUserService(userRepository, transactionRepository, emailService, verificationService, authService)

	return &UserController{
		UserService:         userService,
		EmailService:        emailService,
		VerificationService: verificationService,
		AuthService:         authService,
		PasswordService:     passwordService,
	}
}
