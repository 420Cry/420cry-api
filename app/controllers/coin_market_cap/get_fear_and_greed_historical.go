// Package controllers handles incoming HTTP requests, orchestrates business logic
// through services and repositories, and returns appropriate HTTP responses.
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetFearAndGreedHistorical retrieves the Fear and Greed index
func (h *CoinMarketCapController) GetFearAndGreedHistorical(c *gin.Context) {
	data, err := h.CoinMarketCapService.GetFearAndGreedHistorical()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Success
	c.JSON(http.StatusOK, gin.H{
		"fear_and_greed_historical": data,
	})
}
