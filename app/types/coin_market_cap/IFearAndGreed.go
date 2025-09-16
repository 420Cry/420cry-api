// Package types provides type definitions for FearGreedData
package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// FearGreedData represents the response with a single data object.
type FearGreedData struct {
	Data FearGreedEntry `json:"data"`
}

// FearGreedEntry represents a single fear & greed data point.
type FearGreedEntry struct {
	Value               int       `json:"value"` // Removed the trailing space
	ValueClassification string    `json:"value_classification"`
	UpdateTime          time.Time `json:"update_time"`
}

// FlexibleInt can handle both string and int values during JSON unmarshal.
type FlexibleInt int

// UnmarshalJSON implements the json.Unmarshaler interface for FlexibleInt.
// It attempts to unmarshal the JSON data into an integer, supporting both
// numeric and string representations of integers. If unmarshaling fails,
// it returns an error indicating the failure.
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
