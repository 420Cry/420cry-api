// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"cry-api/app/container"
	walletExplorerService "cry-api/app/services/wallet_explorer"
	EnvTypes "cry-api/app/types/env"
)

// WalletExplorerController handles wallet explorer related requests.
type WalletExplorerController struct {
	TransactionService walletExplorerService.TransactionServiceInterface
}

// NewWalletExplorer initializes a new WalletExplorerController with dependencies from the container.
func NewWalletExplorer(container *container.Container) *WalletExplorerController {
	// Create transaction service (this could be moved to container if needed)
	cfg := container.Get("config").(*EnvTypes.EnvConfig)

	return &WalletExplorerController{
		TransactionService: walletExplorerService.NewTransactionService(cfg),
	}
}
