// Package utils provides utility functions for input sanitization and other helpers.
package utils

import (
	"os"
	"path/filepath"

	"cry-api/app/config"
)

// GenerateEmailTemplatePrefix generates email directory prefix based on environment
func GenerateEmailTemplatePrefix() string {
	// Check if config is loaded, if not use default
	cfg := config.Get()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	var templatePrefix string

	if cfg != nil && cfg.AppEnv == "production" {
		templatePrefix = "/app/app/email/templates"
	} else {
		// Use absolute path for tests and development
		templatePrefix = filepath.Join(cwd, "app", "email", "templates")
	}

	return templatePrefix
}
