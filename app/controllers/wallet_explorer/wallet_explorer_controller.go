package controllers

import (
	EnvTypes "cry-api/app/types/env"
)

// WalletExplorerController handles wallet explorer requests
type WalletExplorerController struct {
	Cfg *EnvTypes.EnvConfig
}

// NewWalletExplorer initializes and returns a WalletExplorerController instance
func NewWalletExplorer(cfg *EnvTypes.EnvConfig) *WalletExplorerController {
	return &WalletExplorerController{
		Cfg: cfg,
	}
}
