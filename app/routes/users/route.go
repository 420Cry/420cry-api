// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/config"
	UserController "cry-api/app/controllers/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the user-related HTTP endpoints to the given Gin router group.
// It initializes the user controller with the database and configuration dependencies.
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Load application config
	cfg := config.Get()

	// Initialize user controller with database and config
	userController := UserController.NewUserController(db, cfg)

	// Route for user signup
	rg.POST("/signup", userController.Signup)

	// Route for verifying a user using the email token (OTP)
	rg.POST("/verify-email-token", userController.VerifyEmailToken)

	// Route for verifying a user using the account token (URL token)
	rg.POST("/verify-account-token", userController.VerifyAccountToken)

	// Route for user signin (login)
	rg.POST("/signin", userController.SignIn)

	// Route for reset password
	rg.POST("/reset-password", userController.HandleResetPasswordRequest)

	// Route for verifying reset password token to save new password
	rg.POST("/verify-reset-password-token", userController.VerifyResetPasswordToken)
}
