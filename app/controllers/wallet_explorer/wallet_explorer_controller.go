// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"cry-api/app/container"
	walletExplorerService "cry-api/app/services/wallet_explorer"
)

// WalletExplorerController handles wallet explorer related requests.
type WalletExplorerController struct {
	TransactionService walletExplorerService.TransactionServiceInterface
}

// NewWalletExplorer initializes a new WalletExplorerController with dependencies from the container.
func NewWalletExplorer(container *container.Container) *WalletExplorerController {
	return &WalletExplorerController{
		TransactionService: container.GetTransactionService(),
	}
}
