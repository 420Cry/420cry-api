// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	WalletExplorerController "cry-api/app/controllers/wallet_explorer"
	"cry-api/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the wallet explorer routes.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	walletExplorerController := WalletExplorerController.NewWalletExplorer(container)

	// Use JWT middleware on this group
	authGroup := rg.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())

	// Authenticated route
	authGroup.GET("/tx", walletExplorerController.GetTransactionInfo)
	authGroup.GET("/xpub", walletExplorerController.GetTransactionByXPUB)
}
