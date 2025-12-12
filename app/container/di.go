// Package container provides improved dependency injection with type safety
package container

import (
	"cry-api/app/config"
	Email "cry-api/app/email"
	UserRepository "cry-api/app/repositories"
	TwoFactorService "cry-api/app/services/2fa"
	AuthService "cry-api/app/services/auth"
	PasswordService "cry-api/app/services/auth/password"
	CoinMarketCapService "cry-api/app/services/coin_market_cap"
	EmailService "cry-api/app/services/email"
	UserService "cry-api/app/services/users"
	WalletExplorerService "cry-api/app/services/wallet_explorer"
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// ServiceContainer is an improved type-safe dependency injection container
type ServiceContainer struct {
	// Core dependencies
	db     *gorm.DB
	config *EnvTypes.EnvConfig

	// Repositories
	userRepo      UserRepository.UserRepository
	userTokenRepo UserRepository.UserTokenRepository

	// Services
	passwordService      PasswordService.PasswordServiceInterface
	emailService         EmailService.EmailServiceInterface
	authService          AuthService.AuthServiceInterface
	userTokenService     UserService.UserTokenServiceInterface
	userService          UserService.UserServiceInterface
	twoFactorService     TwoFactorService.TwoFactorServiceInterface
	coinMarketCapService CoinMarketCapService.CoinMarketCapServiceInterface
	transactionService   WalletExplorerService.TransactionServiceInterface
}

// NewServiceContainer creates a new service container with all dependencies initialized
// This method uses direct initialization. For provider-based initialization, use NewServiceContainerWithProviders
func NewServiceContainer(cfg *EnvTypes.EnvConfig, db *gorm.DB) *ServiceContainer {
	container := &ServiceContainer{
		db:     db,
		config: cfg,
	}

	// Set config globally for backward compatibility
	config.Set(cfg)

	// Initialize repositories
	container.userRepo = UserRepository.NewGormUserRepository(db)
	container.userTokenRepo = UserRepository.NewGormUserTokenRepository(db)

	// Initialize services in dependency order
	container.passwordService = PasswordService.NewPasswordService()

	emailSender := Email.NewSMTPEmailSender(cfg.SMTPConfig.Host, cfg.SMTPConfig.Port)
	emailCreator := &EmailService.EmailCreatorImpl{}
	container.emailService = EmailService.NewEmailService(emailSender, emailCreator)

	container.authService = AuthService.NewAuthService(
		container.userRepo,
		container.passwordService,
	)

	container.userTokenService = UserService.NewUserTokenService(container.userTokenRepo)

	container.userService = UserService.NewUserService(
		container.userRepo,
		container.userTokenRepo,
		container.emailService,
		container.authService,
	)

	container.twoFactorService = TwoFactorService.NewTwoFactorService()
	container.coinMarketCapService = CoinMarketCapService.NewCoinMarketCapServiceService(cfg)
	container.transactionService = WalletExplorerService.NewTransactionService(cfg)

	return container
}

// GetDB returns the database connection
func (c *ServiceContainer) GetDB() *gorm.DB {
	return c.db
}

// GetConfig returns the configuration
func (c *ServiceContainer) GetConfig() *EnvTypes.EnvConfig {
	return c.config
}

// GetUserRepository returns the user repository
func (c *ServiceContainer) GetUserRepository() UserRepository.UserRepository {
	return c.userRepo
}

// GetUserTokenRepository returns the user token repository
func (c *ServiceContainer) GetUserTokenRepository() UserRepository.UserTokenRepository {
	return c.userTokenRepo
}

// GetPasswordService returns the password service
func (c *ServiceContainer) GetPasswordService() PasswordService.PasswordServiceInterface {
	return c.passwordService
}

// GetEmailService returns the email service
func (c *ServiceContainer) GetEmailService() EmailService.EmailServiceInterface {
	return c.emailService
}

// GetAuthService returns the auth service
func (c *ServiceContainer) GetAuthService() AuthService.AuthServiceInterface {
	return c.authService
}

// GetUserTokenService returns the user token service
func (c *ServiceContainer) GetUserTokenService() UserService.UserTokenServiceInterface {
	return c.userTokenService
}

// GetUserService returns the user service
func (c *ServiceContainer) GetUserService() UserService.UserServiceInterface {
	return c.userService
}

// GetTwoFactorService returns the 2FA service
func (c *ServiceContainer) GetTwoFactorService() TwoFactorService.TwoFactorServiceInterface {
	return c.twoFactorService
}

// GetCoinMarketCapService returns the CoinMarketCap service
func (c *ServiceContainer) GetCoinMarketCapService() CoinMarketCapService.CoinMarketCapServiceInterface {
	return c.coinMarketCapService
}

// GetTransactionService returns the transaction/wallet explorer service
func (c *ServiceContainer) GetTransactionService() WalletExplorerService.TransactionServiceInterface {
	return c.transactionService
}
