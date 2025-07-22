// Package controllers handles HTTP requests and responses,
package controllers

import (
	walletExplorerService "cry-api/app/services/wallet_explorer"
	EnvTypes "cry-api/app/types/env"
)

// WalletExplorerController handles wallet explorer related requests.
type WalletExplorerController struct {
	Cfg             *EnvTypes.EnvConfig
	ExternalService walletExplorerService.ExternalServiceInterface
}

// NewWalletExplorer initializes a new WalletExplorerController with the given configuration.
func NewWalletExplorer(cfg *EnvTypes.EnvConfig) *WalletExplorerController {
	return &WalletExplorerController{
		ExternalService: walletExplorerService.NewExternalService(cfg),
	}
}
