package utils

import "github.com/microcosm-cc/bluemonday"

func SanitizeInput(input string) string {
	policy := bluemonday.UGCPolicy()
	return policy.Sanitize(input)
}
