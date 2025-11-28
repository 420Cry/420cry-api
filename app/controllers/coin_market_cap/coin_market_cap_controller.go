// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"cry-api/app/container"
	coinMarketCapService "cry-api/app/services/coin_market_cap"
	EnvTypes "cry-api/app/types/env"
)

// CoinMarketCapController handles coin market cap related requests.
type CoinMarketCapController struct {
	Cfg                  *EnvTypes.EnvConfig
	CoinMarketCapService coinMarketCapService.CoinMarketCapServiceInterface
}

// NewCoinMarketCapController initializes a new CoinMarketCapController with dependencies from the container.
func NewCoinMarketCapController(container *container.Container) *CoinMarketCapController {
	// Get config from container
	cfg := container.Get("config").(*EnvTypes.EnvConfig)

	return &CoinMarketCapController{
		Cfg:                  cfg,
		CoinMarketCapService: coinMarketCapService.NewCoinMarketCapServiceService(cfg),
	}
}
