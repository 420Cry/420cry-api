// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/config"
	controller "cry-api/app/controllers/2fa"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the user-related HTTP endpoints to the given Gin router group.
// It initializes the user controller with the database and configuration dependencies.
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Load application config
	cfg := config.Get()

	// Initialize user controller with database and config
	TwoFactorController := controller.NewTwoFactorController(db, cfg)

	// Route for user setup
	rg.POST("/setup", TwoFactorController.Setup)

	// Route for verify 2fa otp
	rg.POST("/verify-otp", TwoFactorController.VerifyOTP)
}
