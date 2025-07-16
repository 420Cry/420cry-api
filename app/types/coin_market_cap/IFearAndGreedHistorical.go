// Package types provides type definitions for FearGreedData
package types

// FearGreedHistorical represents the response with a single data object.
type FearGreedHistorical struct {
	Data []FearGreedDataPoint `json:"data"`
}

// FearGreedDataPoint represents a single historical data point of the Fear and Greed index.
type FearGreedDataPoint struct {
	Timestamp           string `json:"timestamp"`            // ISO date string or UNIX timestamp as string
	Value               int    `json:"value"`                // Fear and Greed index value
	ValueClassification string `json:"value_classification"` // Classification like "Fear", "Extreme Fear", etc.
}
