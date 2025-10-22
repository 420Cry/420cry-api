// Package types defines environment configuration types for the application.
package types

import (
	"errors"
	"fmt"
)

// SMTPConfig holds the configuration for SMTP settings
type SMTPConfig struct {
	Host string
	Port string
}

// WalletExplorerConfig holds external API configuration for wallet explorer services.
type WalletExplorerConfig struct {
	API string
}

// BlockchainConfig holds external API configuration for wallet explorer services.
type BlockchainConfig struct {
	API string
}

// CoinMarketCapConfig holds external API configuration for coin market cap services.
type CoinMarketCapConfig struct {
	API    string
	APIKey string
}

// EnvConfig maps environment variables to application configuration fields.
type EnvConfig struct {
	AppEnv               string
	CryAppURL            string
	CryAPIURL            string
	APIPort              int
	DBHost               string
	DBPort               int
	DBDatabase           string
	DBUserName           string
	DBPassword           string
	SMTPConfig           SMTPConfig
	NoReplyEmail         string
	WalletExplorerConfig WalletExplorerConfig
	BlockchainConfig     BlockchainConfig
	CoinMarketCapConfig  CoinMarketCapConfig
}

// Validate validates the configuration
func (c *EnvConfig) Validate() error {
	// Validate required fields
	if c.DBHost == "" {
		return errors.New("DB_HOST is required")
	}

	if c.DBDatabase == "" {
		return errors.New("DB_DATABASE is required")
	}

	if c.DBUserName == "" {
		return errors.New("DB_USERNAME is required")
	}

	if c.DBPassword == "" {
		return errors.New("DB_PASSWORD is required")
	}

	if c.APIPort <= 0 || c.APIPort > 65535 {
		return fmt.Errorf("API_PORT must be between 1 and 65535, got %d", c.APIPort)
	}

	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535, got %d", c.DBPort)
	}

	if c.NoReplyEmail == "" {
		return errors.New("NO_REPLY_EMAIL is required")
	}

	if c.CryAppURL == "" {
		return errors.New("CRY_APP_URL is required")
	}

	if c.CryAPIURL == "" {
		return errors.New("CRY_API_URL is required")
	}

	// Validate SMTP configuration
	if c.SMTPConfig.Host == "" {
		return errors.New("SMTP_HOST is required")
	}

	if c.SMTPConfig.Port == "" {
		return errors.New("SMTP_PORT is required")
	}

	return nil
}
