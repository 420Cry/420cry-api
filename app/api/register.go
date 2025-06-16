// Package api provides route registration and API endpoint setup for the application.
package api

import (
	users "cry-api/app/api/routes/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAllRoutes sets up all API routes using Gin.
func RegisterAllRoutes(r *gin.Engine, db *gorm.DB) {
	users.RegisterRoutes(r.Group("/users"), db)
}
