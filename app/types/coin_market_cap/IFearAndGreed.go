// Package types provides type definitions for user signup requests.
package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// FearGreedData represents the response with a single data object.
type FearGreedData struct {
	Data   FearGreedEntry `json:"data"`
	Status Status         `json:"status"`
}

// FearGreedEntry represents a single fear & greed data point.
type FearGreedEntry struct {
	Value               int       `json:"value"` // Removed the trailing space
	ValueClassification string    `json:"value_classification"`
	UpdateTime          time.Time `json:"update_time"`
}

// Status represents metadata about the API response.
type Status struct {
	Timestamp    time.Time   `json:"timestamp"`
	ErrorCode    FlexibleInt `json:"error_code"` // Custom type to handle string or int
	ErrorMessage string      `json:"error_message"`
	Elapsed      int         `json:"elapsed"`
	CreditCount  int         `json:"credit_count"`
	Notice       string      `json:"notice"`
}

// FlexibleInt can handle both string and int values during JSON unmarshal.
type FlexibleInt int

func (f *FlexibleInt) UnmarshalJSON(data []byte) error {
	// Try as int
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*f = FlexibleInt(intVal)
		return nil
	}

	// Try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		_, err := fmt.Sscanf(strVal, "%d", &intVal)
		if err == nil {
			*f = FlexibleInt(intVal)
			return nil
		}
	}

	return fmt.Errorf("cannot unmarshal error_code into FlexibleInt: %s", string(data))
}
