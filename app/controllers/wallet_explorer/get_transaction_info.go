// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"

	services "cry-api/app/services/jwt"

	"github.com/gin-gonic/gin"
)

// GetTransactionInfo retrieves transaction information by transaction ID.
func (h *WalletExplorerController) GetTransactionInfo(c *gin.Context) {
	// User validation
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User info not found in context"})
		return
	}
	_ = claims.(*services.Claims)

	// Query param
	txid := c.Query("txid")
	if txid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing txid parameter"})
		return
	}

	// Call external service
	data, err := h.ExternalService.GetTransactionByTxID(txid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"isSuccess": false,
			"message":   "Failed to retrieve transaction data",
			"error":     err.Error(),
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"transaction_data": data,
	})
}
