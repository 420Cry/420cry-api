// Package routes provides route registration and API endpoint setup for the application.
package routes

import (
	TwoFactorRoute "cry-api/app/routes/2fa"
	OAuthRoute "cry-api/app/routes/oauth"
	UserRoute "cry-api/app/routes/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAllRoutes sets up all API routes using Gin.
func RegisterAllRoutes(r *gin.Engine, db *gorm.DB) {
	UserRoute.RegisterRoutes(r.Group("/users"), db)
	TwoFactorRoute.RegisterRoutes(r.Group("/2fa"), db)
	OAuthRoute.RegisterRoutes(r.Group("/auth"), db)
}
