// Package services provides  wallet explorer services for external API interactions.
package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	EnvTypes "cry-api/app/types/env"
	WalletExplorer "cry-api/app/types/wallet_explorer"
)

// ExternalService interacts with external wallet explorer APIs.
type ExternalService struct {
	Config *EnvTypes.EnvConfig
}

// NewExternalService initializes and returns an ExternalService instance
func NewExternalService(cfg *EnvTypes.EnvConfig) *ExternalService {
	return &ExternalService{
		Config: cfg,
	}
}

// GetTransactionByTxID fetches transaction data from WalletExplorer API
func (s *ExternalService) GetTransactionByTxID(txid string) (*WalletExplorer.ITransactionData, error) {
	// Use config URL
	baseURL := s.Config.BlockchainConfig.API
	url := fmt.Sprintf("%s/rawtx/%s", baseURL, txid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data WalletExplorer.ITransactionData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &data, nil
}
