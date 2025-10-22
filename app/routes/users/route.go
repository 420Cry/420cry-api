// Package routes sets up the HTTP routing for the application.
package routes

import (
	"cry-api/app/container"
	UserController "cry-api/app/controllers/users"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the user-related HTTP endpoints to the given Gin router group.
// It initializes the user controller with dependencies from the container.
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	// Initialize user controller with container dependencies
	userController := UserController.NewUserController(container)

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
