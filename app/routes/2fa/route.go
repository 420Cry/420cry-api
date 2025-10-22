// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	controller "cry-api/app/controllers/2fa"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the 2FA-related HTTP endpoints to the given Gin router group.
// It initializes the 2FA controller with dependencies from the container.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	// Initialize 2FA controller with container dependencies
	TwoFactorController := controller.NewTwoFactorController(container)

	// Route for user setup
	rg.POST("/setup", TwoFactorController.Setup)

	// Route for verify 2fa setup - otp
	rg.POST("/setup/verify-otp", TwoFactorController.VerifySetUpOTP)

	// Route for verify 2fa (AUTH)
	rg.POST("/auth/verify-otp", TwoFactorController.VerifyOTP)

	// Route for alternative - send otp to user email
	rg.POST("/alternative/send-email-otp", TwoFactorController.AlternativeSendOtp)

	// Route for verify alternative otp from user email
	rg.POST("/alternative/verify-email-otp", TwoFactorController.AlternativeVerifyOTP)
}
