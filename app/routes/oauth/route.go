// Package routes provides route registration and API endpoint setup for the application.
package routes

import (
	"cry-api/app/container"
	OAuthController "cry-api/app/controllers/oauth"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the OAuth related HTTP endpoints to the given Gin router group
// It initializes controllers to be used for handle specific endpoints
func RegisterRoutes(rg *gin.RouterGroup, container *container.Container) {
	OAuthController := OAuthController.NewOAuthController(container)

	rg.GET("/google/callback", OAuthController.HandleGoogleCallback)
	rg.GET("/discord/callback")
}
