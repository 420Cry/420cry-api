// Package routes provides route registration and API endpoint setup for the application.
package routes

import (
	"cry-api/app/config"
	OAuthController "cry-api/app/controllers/oauth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the OAuth related HTTP endpoints to the given Gin router group
// It initializes controllers to be used for handle specific endpoints
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	cfg := config.Get()

	OAuthController := OAuthController.NewOAuthController(db, cfg)

	rg.GET("/google/callback", OAuthController.HandleGoogleCallback)
	rg.GET("/discord/callback")
}
