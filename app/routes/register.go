// Package routes provides route registration and API endpoint setup for the application.
package routes

import (
	"cry-api/app/container"
	TwoFactorRoute "cry-api/app/routes/2fa"
	CoinMarketRoute "cry-api/app/routes/coin_market_cap"
	OAuthRoute "cry-api/app/routes/oauth"
	UserRoute "cry-api/app/routes/users"
	WalletExplorerRoute "cry-api/app/routes/wallet_explorer"

	"github.com/gin-gonic/gin"
)

// RegisterAllRoutes sets up all API routes using Gin with dependency injection container.
func RegisterAllRoutes(r *gin.Engine, container *container.Container) {
	// API versioning
	v1 := r.Group("/api/v1")

	// Seperate routes for OAuth
	oauthRoute := r.Group("")

	UserRoute.RegisterRoutes(v1.Group("/users"), container)
	TwoFactorRoute.RegisterRoutes(v1.Group("/2fa"), container)
	OAuthRoute.RegisterRoutes(oauthRoute.Group("/auth"), container)
	WalletExplorerRoute.RegisterRoutes(v1.Group("/wallet-explorer"), container)
	CoinMarketRoute.RegisterRoutes(v1.Group("/coin-market-cap"), container)
}
