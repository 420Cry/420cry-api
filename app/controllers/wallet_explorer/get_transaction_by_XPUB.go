// Package controllers handles HTTP requests and responses,
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTransactionByXPUB retrieves transactions by transaction XPUB.
func (h *WalletExplorerController) GetTransactionByXPUB(c *gin.Context) {
	xpub := c.Query("xpub")
	if xpub == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing xpub parameter"})
		return
	}

	data, err := h.TransactionService.GetTransactionByXPUB(xpub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"xpub": data})
}
