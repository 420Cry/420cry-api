// Package services provides  wallet explorer services for external API interactions.
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	CoinMarketCap "cry-api/app/types/coin_market_cap"
	EnvTypes "cry-api/app/types/env"
)

// CoinMarketCapService interacts with external wallet explorer APIs.
type CoinMarketCapService struct {
	Config *EnvTypes.EnvConfig
}

// CoinMarketCapServiceInterface defines the methods for the CoinMarketCapService.
type CoinMarketCapServiceInterface interface {
	GetFearAndGreed() (*CoinMarketCap.FearGreedData, error)
}

// NewCoinMarketCapServiceService initializes and returns an CoinMarketCapService instance
func NewCoinMarketCapServiceService(cfg *EnvTypes.EnvConfig) *CoinMarketCapService {
	return &CoinMarketCapService{
		Config: cfg,
	}
}

// GetFearAndGreed fetches fear and greed index data from coin market cap API
func (s *CoinMarketCapService) GetFearAndGreed() (*CoinMarketCap.FearGreedData, error) {
	// Use config URL
	baseURL := s.Config.CoinMarketCapConfig.API
	url := fmt.Sprintf("%s/v3/fear-and-greed/historical", baseURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Optional: Add headers like API key if required by the API
	if s.Config.CoinMarketCapConfig.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Config.CoinMarketCapConfig.APIKey))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Check for non-200 response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Decode response body into struct
	var data CoinMarketCap.FearGreedData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &data, nil
}
