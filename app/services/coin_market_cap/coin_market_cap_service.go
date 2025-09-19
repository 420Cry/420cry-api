// Package services provides  coin market cap services for external API interactions.
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
	GetFearAndGreedLastest() (*CoinMarketCap.FearGreedData, error)
	GetFearAndGreedHistorical(start, limit int) (*CoinMarketCap.FearGreedHistorical, error)
}

// NewCoinMarketCapServiceService initializes and returns an CoinMarketCapService instance
func NewCoinMarketCapServiceService(cfg *EnvTypes.EnvConfig) *CoinMarketCapService {
	return &CoinMarketCapService{
		Config: cfg,
	}
}

// GetFearAndGreedLastest fetches fear and greed index data from CoinMarketCap API.
func (s *CoinMarketCapService) GetFearAndGreedLastest() (*CoinMarketCap.FearGreedData, error) {
	baseURL := s.Config.CoinMarketCapConfig.API
	url := fmt.Sprintf("%s/v3/fear-and-greed/latest", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Use the correct custom header for API key authentication
	req.Header.Set("X-CMC_PRO_API_KEY", s.Config.CoinMarketCapConfig.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var data CoinMarketCap.FearGreedData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &data, nil
}

// GetFearAndGreedHistorical fetches fear and greed index data from CoinMarketCap API.
func (s *CoinMarketCapService) GetFearAndGreedHistorical(start, limit int) (*CoinMarketCap.FearGreedHistorical, error) {
	// Ensure the limit does not exceed 500
	if limit > 500 {
		limit = 500
	}
	if limit < 1 {
		limit = 1
	}
	if start < 1 {
		start = 1
	}

	baseURL := s.Config.CoinMarketCapConfig.API
	url := fmt.Sprintf("%s/v3/fear-and-greed/historical?start=%d&limit=%d", baseURL, start, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("X-CMC_PRO_API_KEY", s.Config.CoinMarketCapConfig.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var data CoinMarketCap.FearGreedHistorical
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &data, nil
}
