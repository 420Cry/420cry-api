// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	CoinMarketCapController "cry-api/app/controllers/coin_market_cap"
	"cry-api/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the coin market cap routes.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	coinMarketCapController := CoinMarketCapController.NewCoinMarketCapController(container)

	// Use JWT middleware on this group
	authGroup := rg.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())

	// Authenticated route
	authGroup.GET("/fear-and-greed-lastest", coinMarketCapController.GetFearAndGreedLastest)
	authGroup.GET("/fear-and-greed-historical", coinMarketCapController.GetFearAndGreedHistorical)
}
