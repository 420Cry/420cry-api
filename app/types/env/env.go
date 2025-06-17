// Package types defines environment configuration types for the application.
package types

// SMTPConfig holds the configuration for SMTP settings
type SMTPConfig struct {
	Host string
	Port string
}

// EnvConfig maps environment variables to application configuration fields.
type EnvConfig struct {
	CryAppURL    string
	CryApiURL		 string
	APIPort      int
	DBHost       string
	DBPort       int
	DBDatabase   string
	DBUserName   string
	DBPassword   string
	SMTPConfig   SMTPConfig
	NoReplyEmail string
}
