// Package types provides type definitions for user signup requests.
package types

import "time"

// FearGreedData represents the structure of the Fear and Greed index data.
type FearGreedData struct {
	Data   []FearGreedEntry `json:"data"`
	Status Status           `json:"status"`
}

// FearGreedEntry represents a single entry in the Fear and Greed index.
type FearGreedEntry struct {
	Timestamp           string `json:"timestamp"`            // Unix timestamp as string
	Value               int    `json:"value"`                // e.g., 38
	ValueClassification string `json:"value_classification"` // e.g., "Fear"
}

// Status represents the status of the API response.
type Status struct {
	Timestamp    time.Time `json:"timestamp"` // ISO 8601 format
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
	Notice       string    `json:"notice"`
}
