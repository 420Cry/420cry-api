// Package utils provides utility functions for input sanitization and other helpers.
package utils

import "github.com/microcosm-cc/bluemonday"

// SanitizeInput removes potentially dangerous HTML from the input string
// using a UGC (User Generated Content) policy to prevent XSS attacks.
func SanitizeInput(input string) string {
	policy := bluemonday.UGCPolicy()
	return policy.Sanitize(input)
}
