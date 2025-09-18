// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	coinMarketCapService "cry-api/app/services/coin_market_cap"
	EnvTypes "cry-api/app/types/env"
)

// CoinMarketCapController handles wallet explorer related requests.
type CoinMarketCapController struct {
	Cfg                  *EnvTypes.EnvConfig
	CoinMarketCapService coinMarketCapService.CoinMarketCapServiceInterface
}

// NewCoinMarketCapController initializes a new CoinMarketCapController with the given configuration.
func NewCoinMarketCapController(cfg *EnvTypes.EnvConfig) *CoinMarketCapController {
	return &CoinMarketCapController{
		CoinMarketCapService: coinMarketCapService.NewCoinMarketCapServiceService(cfg),
	}
}
