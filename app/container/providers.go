// Package container provides service providers for organizing microservices
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

// ServiceProvider defines the interface for service providers
type ServiceProvider interface {
	Register(container *ServiceContainer)
}

// UserServiceProvider registers user-related services
type UserServiceProvider struct{}

// Register initializes user repositories and token services
func (p *UserServiceProvider) Register(c *ServiceContainer) {
	c.userRepo = UserRepository.NewGormUserRepository(c.db)
	c.userTokenRepo = UserRepository.NewGormUserTokenRepository(c.db)
	c.userTokenService = UserService.NewUserTokenService(c.userTokenRepo)
}

// AuthServiceProvider registers authentication-related services
type AuthServiceProvider struct{}

// Register initializes password and authentication services
func (p *AuthServiceProvider) Register(c *ServiceContainer) {
	c.passwordService = PasswordService.NewPasswordService()
	c.authService = AuthService.NewAuthService(
		c.userRepo,
		c.passwordService,
	)
}

// EmailServiceProvider registers email-related services
type EmailServiceProvider struct{}

// Register initializes email sender and creator services
func (p *EmailServiceProvider) Register(c *ServiceContainer) {
	emailSender := Email.NewSMTPEmailSender(c.config.SMTPConfig.Host, c.config.SMTPConfig.Port)
	emailCreator := &EmailService.EmailCreatorImpl{}
	c.emailService = EmailService.NewEmailService(emailSender, emailCreator)
}

// UserBusinessServiceProvider registers user business logic services
type UserBusinessServiceProvider struct{}

// Register initializes user business service with all dependencies
func (p *UserBusinessServiceProvider) Register(c *ServiceContainer) {
	c.userService = UserService.NewUserService(
		c.userRepo,
		c.userTokenRepo,
		c.emailService,
		c.authService,
	)
}

// TwoFactorServiceProvider registers 2FA services
type TwoFactorServiceProvider struct{}

// Register initializes two-factor authentication service
func (p *TwoFactorServiceProvider) Register(c *ServiceContainer) {
	c.twoFactorService = TwoFactorService.NewTwoFactorService()
}

// ExternalAPIServiceProvider registers external API services
type ExternalAPIServiceProvider struct{}

// Register initializes external API services (CoinMarketCap and Wallet Explorer)
func (p *ExternalAPIServiceProvider) Register(c *ServiceContainer) {
	c.coinMarketCapService = CoinMarketCapService.NewCoinMarketCapServiceService(c.config)
	c.transactionService = WalletExplorerService.NewTransactionService(c.config)
}

// registerAllProviders registers all service providers in the correct order
func registerAllProviders(container *ServiceContainer) {
	providers := []ServiceProvider{
		&UserServiceProvider{},
		&AuthServiceProvider{},
		&EmailServiceProvider{},
		&UserBusinessServiceProvider{},
		&TwoFactorServiceProvider{},
		&ExternalAPIServiceProvider{},
	}

	for _, provider := range providers {
		provider.Register(container)
	}
}

// NewServiceContainerWithProviders creates a service container using providers
// This is an alternative initialization method that uses the provider pattern
func NewServiceContainerWithProviders(cfg *EnvTypes.EnvConfig, db *gorm.DB) *ServiceContainer {
	container := &ServiceContainer{
		db:     db,
		config: cfg,
	}

	// Set config globally for backward compatibility
	config.Set(cfg)

	// Register all services using providers
	registerAllProviders(container)

	return container
}
