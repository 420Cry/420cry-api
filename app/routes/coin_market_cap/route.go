// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	CoinMarketCapController "cry-api/app/controllers/coin_market_cap"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the coin market cap routes.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	coinMarketCapController := CoinMarketCapController.NewCoinMarketCapController(container)

	// Public routes (no authentication required)
	rg.GET("/fear-and-greed-lastest", coinMarketCapController.GetFearAndGreedLastest)
	rg.GET("/fear-and-greed-historical", coinMarketCapController.GetFearAndGreedHistorical)
}
