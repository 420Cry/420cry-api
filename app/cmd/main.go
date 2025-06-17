// Package main provides server configuration and start up logic
package main

import (
	"log"
	"strconv"
	"time"

	"cry-api/app/config"
	"cry-api/app/database"
	"cry-api/app/routes"
	Env "cry-api/app/types/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Get()

	dbConn, err := database.GetDBConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	db := dbConn.GetDB()

	router := gin.Default()

	router.Use(SetupCORS(cfg))

	routes.RegisterAllRoutes(router, db)

	if err := router.Run(":" + strconv.Itoa(cfg.APIPort)); err != nil {
		log.Fatal("Failed to run server: ", err)
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
