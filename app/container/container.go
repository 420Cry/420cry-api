// Package container provides dependency injection functionality for the application.
package container

import (
	"cry-api/app/config"
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	AuthService "cry-api/app/services/auth"
	PasswordService "cry-api/app/services/auth/password"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// Container manages application dependencies
type Container struct {
	services map[string]interface{}
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

// Register adds a service to the container
func (c *Container) Register(name string, service interface{}) {
	c.services[name] = service
}

// Get retrieves a service from the container
func (c *Container) Get(name string) interface{} {
	return c.services[name]
}

// InitializeContainer sets up all application dependencies
func InitializeContainer(cfg *EnvTypes.EnvConfig, db *gorm.DB) *Container {
	container := NewContainer()

	// Register database and config
	container.Register("database", db)
	container.Register("config", cfg)

	// Set the config globally for any packages that need it
	config.Set(cfg)

	// Register repositories
	userRepo := UserRepository.NewGormUserRepository(db)
	userTokenRepo := UserRepository.NewGormUserTokenRepository(db)
	container.Register("userRepository", userRepo)
	container.Register("userTokenRepository", userTokenRepo)

	// Register services
	passwordService := PasswordService.NewPasswordService()
	container.Register("passwordService", passwordService)

	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailCreator := &EmailService.EmailCreatorImpl{}
	emailService := EmailService.NewEmailService(emailSender, emailCreator)
	container.Register("emailService", emailService)

	authService := AuthService.NewAuthService(userRepo, passwordService)
	container.Register("authService", authService)

	userTokenService := UserService.NewUserTokenService(userTokenRepo)
	container.Register("userTokenService", userTokenService)

	userService := UserService.NewUserService(
		userRepo,
		userTokenRepo,
		emailService,
		authService,
	)
	container.Register("userService", userService)

	return container
}

// GetUserRepository returns the user repository from container
func (c *Container) GetUserRepository() UserRepository.UserRepository {
	return c.Get("userRepository").(UserRepository.UserRepository)
}

// GetUserTokenRepository returns the user token repository from container
func (c *Container) GetUserTokenRepository() UserRepository.UserTokenRepository {
	return c.Get("userTokenRepository").(UserRepository.UserTokenRepository)
}

// GetUserService returns the user service from container
func (c *Container) GetUserService() UserService.UserServiceInterface {
	return c.Get("userService").(UserService.UserServiceInterface)
}

// GetUserTokenService returns the user token service from container
func (c *Container) GetUserTokenService() UserService.UserTokenServiceInterface {
	return c.Get("userTokenService").(UserService.UserTokenServiceInterface)
}

// GetEmailService returns the email service from container
func (c *Container) GetEmailService() EmailService.EmailServiceInterface {
	return c.Get("emailService").(EmailService.EmailServiceInterface)
}

// GetAuthService returns the auth service from container
func (c *Container) GetAuthService() AuthService.AuthServiceInterface {
	return c.Get("authService").(AuthService.AuthServiceInterface)
}

// GetPasswordService returns the password service from container
func (c *Container) GetPasswordService() PasswordService.PasswordServiceInterface {
	return c.Get("passwordService").(PasswordService.PasswordServiceInterface)
}
