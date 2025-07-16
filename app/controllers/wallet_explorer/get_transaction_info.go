// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTransactionInfo retrieves transaction information by transaction ID.
func (h *WalletExplorerController) GetTransactionInfo(c *gin.Context) {
	// Query param
	txid := c.Query("txid")
	if txid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing txid parameter"})
		return
	}

	data, err := h.TransactionService.GetTransactionByTxID(txid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"transaction_data": data,
	})
}
