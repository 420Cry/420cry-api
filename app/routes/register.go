// Package routes provides route registration and API endpoint setup for the application.
package routes

import (
	userroute "cry-api/app/routes/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAllRoutes sets up all API routes using Gin.
func RegisterAllRoutes(r *gin.Engine, db *gorm.DB) {
	userroute.RegisterRoutes(r.Group("/users"), db)
}
