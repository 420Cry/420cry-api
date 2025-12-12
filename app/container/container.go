// Package container provides dependency injection functionality for the application.
// This file maintains backward compatibility with the old Container type.
package container

import (
	EnvTypes "cry-api/app/types/env"

	"gorm.io/gorm"
)

// Container is the legacy container interface for backward compatibility
// It wraps the new ServiceContainer to maintain existing API
type Container struct {
	*ServiceContainer
}

// InitializeContainer creates a new container with all dependencies
// This maintains backward compatibility while using the new ServiceContainer internally
func InitializeContainer(cfg *EnvTypes.EnvConfig, db *gorm.DB) *Container {
	return &Container{
		ServiceContainer: NewServiceContainer(cfg, db),
	}
}

// Get is a legacy method for backward compatibility
// It allows accessing services by string name (not recommended for new code)
func (c *Container) Get(name string) interface{} {
	switch name {
	case "database":
		return c.GetDB()
	case "config":
		return c.GetConfig()
	case "userRepository":
		return c.GetUserRepository()
	case "userTokenRepository":
		return c.GetUserTokenRepository()
	case "passwordService":
		return c.GetPasswordService()
	case "emailService":
		return c.GetEmailService()
	case "authService":
		return c.GetAuthService()
	case "userTokenService":
		return c.GetUserTokenService()
	case "userService":
		return c.GetUserService()
	case "twoFactorService":
		return c.GetTwoFactorService()
	case "coinMarketCapService":
		return c.GetCoinMarketCapService()
	case "transactionService":
		return c.GetTransactionService()
	default:
		return nil
	}
}

// Register is a legacy method for backward compatibility (no-op in new implementation)
func (c *Container) Register(_ string, _ interface{}) {
	// No-op: services are initialized in NewServiceContainer
	// This method exists only for backward compatibility
}
