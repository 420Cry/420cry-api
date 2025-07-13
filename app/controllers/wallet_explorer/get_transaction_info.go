// Package controllers handles HTTP requests and responses.
package controllers

import (
	services "cry-api/app/services/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTransactionInfo
func (h *WalletExplorerController) GetTransactionInfo(c *gin.Context) {
	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User info not found in context"})
		return
	}

	userClaims := claims.(*services.Claims)

	txid := c.Query("txid")
	if txid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing txid parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction info retrieved successfully",
		"txid":    txid,
		"user": gin.H{
			"uuid":  userClaims.UUID,
			"email": userClaims.Email,
		},
	})
}
