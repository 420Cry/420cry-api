// Package utils provides utility functions for input sanitization and other helpers.
package utils

import (
	"cry-api/app/config"
)

// GenerateEmailTemplatePrefix generates email directory prefix based on environment
func GenerateEmailTemplatePrefix() string {
	AppEnv := config.Get().AppEnv
	var templatePrefix string

	if AppEnv == "production" {
		templatePrefix = "/app/app/email/templates"
		return templatePrefix
	}

	templatePrefix = "app/email/templates"
	return templatePrefix
}
