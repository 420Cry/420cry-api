// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/config"
	CoinMarketCapController "cry-api/app/controllers/coin_market_cap"
	"cry-api/app/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the wallet explorer routes.
func RegisterRoutes(rg *gin.RouterGroup) {
	cfg := config.Get()

	coinMarketCapController := CoinMarketCapController.NewCoinMarketCapController(cfg)

	// Use JWT middleware on this group
	authGroup := rg.Group("")
	authGroup.Use(middleware.JWTAuthMiddleware())

	// Authenticated route
	authGroup.GET("/fear-and-greed-lastest", coinMarketCapController.GetFearAndGreedLastest)
	authGroup.GET("/fear-and-greed-historical", coinMarketCapController.GetFearAndGreedHistorical)
}
