// Package main provides server configuration and start up logic
package main

import (
	"strconv"
	"time"

	"cry-api/app/config"
	"cry-api/app/container"
	"cry-api/app/database"
	"cry-api/app/logger"
	"cry-api/app/middleware"
	"cry-api/app/routes"
	Env "cry-api/app/types/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	appLogger := logger.GetLogger()
	appLogger.Info("Starting 420cry API server")

	// Load and validate configuration
	cfg := config.Get()
	appLogger.Info("Configuration loaded successfully")

	// Initialize database connection
	dbConn, err := database.GetDBConnection()
	if err != nil {
		appLogger.WithError(err).Fatal("Database connection failed")
	}

	// Configure connection pool
	if err := dbConn.ConfigureConnectionPool(25, 5, time.Hour); err != nil {
		appLogger.WithError(err).Fatal("Failed to configure database connection pool")
	}

	// Test database connection
	if err := dbConn.Ping(); err != nil {
		appLogger.WithError(err).Fatal("Database ping failed")
	}

	db := dbConn.GetDB()
	appLogger.Info("Database connection established successfully")

	// Initialize dependency injection container
	container := container.InitializeContainer(cfg, db)
	appLogger.Info("Dependency injection container initialized")

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.SecurityMiddleware())
	router.Use(middleware.RequestLoggerMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.ContentTypeMiddleware())
	router.Use(middleware.RequestSizeMiddleware(10 * 1024 * 1024)) // 10MB limit
	router.Use(middleware.HealthCheckMiddleware())
	router.Use(SetupCORS(cfg))

	// Register routes with container
	routes.RegisterAllRoutes(router, container)

	appLogger.WithField("port", cfg.APIPort).Info("Server starting on port")

	if err := router.Run(":" + strconv.Itoa(cfg.APIPort)); err != nil {
		appLogger.WithError(err).Fatal("Failed to run server")
	}
}

// SetupCORS funcs provides config and setups CORS
func SetupCORS(cfg *Env.EnvConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CryAppURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,

		AllowOriginFunc: func(origin string) bool {
			return origin == "" || origin == cfg.CryAppURL
		},
	})
}
