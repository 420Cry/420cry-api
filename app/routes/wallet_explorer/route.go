// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/config"
	WalletExplorerController "cry-api/app/controllers/wallet_explorer"
	"cry-api/app/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	cfg := config.Get()

	walletExplorerController := WalletExplorerController.NewWalletExplorer(cfg)

	// Use JWT middleware on this group
	authGroup := rg.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())

	// Authenticated route
	authGroup.GET("/tx", walletExplorerController.GetTransactionInfo)
}
