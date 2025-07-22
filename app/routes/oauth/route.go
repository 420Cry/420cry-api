package oauth

import (
	"cry-api/app/config"
	OAuthController "cry-api/app/controllers/oauth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	cfg := config.Get()

	OAuthController := OAuthController.NewOAuthController(db, cfg)
	
	rg.GET("/google/callback", OAuthController.HandleGoogleCallback)
	rg.GET("/discord/callback")
}
