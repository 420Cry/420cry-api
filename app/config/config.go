// Package config provides functionality for loading and managing application configuration from environment variables.
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	types "cry-api/app/types/env"

	"github.com/joho/godotenv"
)

var (
	configInstance *types.EnvConfig
	configLoaded   = false
)

// Set allows setting the config instance manually
func Set(cfg *types.EnvConfig) {
	configInstance = cfg
	configLoaded = true
}

// Get retrieves the loaded config instance, loading it if necessary
func Get() *types.EnvConfig {
	if !configLoaded {
		Load()
	}
	return configInstance
}

// Load loads environment variables and sets the config instance
func Load() *types.EnvConfig {
	if configLoaded {
		return configInstance
	}

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// AppEnv
	appEnv := os.Getenv("APP_ENV")

	// CryAppUrl
	cryAppURL := os.Getenv("CRY_APP_URL")
	if !strings.HasPrefix(cryAppURL, "http://") && !strings.HasPrefix(cryAppURL, "https://") {
		cryAppURL = "https://" + cryAppURL
	}

	// CryAPIURL
	CryAPIURL := os.Getenv("CRY_API_URL")
	if !strings.HasPrefix(CryAPIURL, "http://") && !strings.HasPrefix(CryAPIURL, "https://") {
		CryAPIURL = "https://" + CryAPIURL
	}

	// Load API Port with a fallback value
	apiPort := getEnvAsInt("API_PORT", 8080)

	// Load DB Port with a fallback value
	dbPort := getEnvAsInt("DB_PORT", 3306)

	dbHost := os.Getenv("DB_HOST")
	db := os.Getenv("DB_DATABASE")
	mysqlUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Load SMTP variables with defaults
	smtpHost := getEnv("SMTP_HOST", "mailhog")
	smtpPort := getEnv("SMTP_PORT", "1025")

	// Load NoReplyEmail
	noReplyEmail := os.Getenv("NO_REPLY_EMAIL")

	// Load WALLET_EXPLORER_API
	walletExplorerApi := os.Getenv("WALLET_EXPLORER_API")

	// Set the config instance
	configInstance = &types.EnvConfig{
		AppEnv:       appEnv,
		CryAppURL:    cryAppURL,
		CryAPIURL:    CryAPIURL,
		APIPort:      apiPort,
		DBHost:       dbHost,
		DBPort:       dbPort,
		DBDatabase:   db,
		DBUserName:   mysqlUser,
		DBPassword:   dbPassword,
		NoReplyEmail: noReplyEmail,
		SMTPConfig: types.SMTPConfig{
			Host: smtpHost,
			Port: smtpPort,
		},
		WalletExplorerConfig: types.WalletExplorerConfig{
			API: walletExplorerApi,
		},
	}

	configLoaded = true
	return configInstance
}

// Reload forces a reload of the config
func Reload() {
	configLoaded = false
	Load()
}

// Helper function to get an environment variable as a string with a fallback value
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// Helper function to get an environment variable as an integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}
