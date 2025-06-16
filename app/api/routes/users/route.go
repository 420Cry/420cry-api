package users

import (
	"cry-api/app/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers the user-related HTTP routes with the provided Gin router group.
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	cfg := config.Get()
	handler := NewHandler(db, cfg)

	rg.POST("/signup", handler.Signup)
	rg.POST("/verify-email-token", handler.VerifyEmailToken)
	rg.POST("/verify-account-token", handler.VerifyAccountToken)
	rg.POST("/signin", handler.SignIn)
}
