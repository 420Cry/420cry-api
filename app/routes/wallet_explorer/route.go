// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	WalletExplorerController "cry-api/app/controllers/wallet_explorer"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the wallet explorer routes.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	walletExplorerController := WalletExplorerController.NewWalletExplorer(container)

	// Public routes (no authentication required)
	rg.GET("/tx", walletExplorerController.GetTransactionInfo)
	rg.GET("/xpub", walletExplorerController.GetTransactionByXPUB)
}
