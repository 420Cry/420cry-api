package main

import (
	"420-api/app/config"
	"420-api/app/middleware"
	"420-api/app/routes"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Load the configuration
	cfg := config.Load()

	// Create a new router
	r := mux.NewRouter()

	// Register routes
	routes.RegisterHelloRoute(r)

	// Initialize middlewares
	limiter := middleware.NewRateLimiter(1, 5)

	// Add Bearer Token Middleware and Rate Limiter
	secureRouter := middleware.SecurityHeaders(
		limiter.Middleware(r),
	)

	// Start the server using the API_PORT from the config
	log.Printf("Server started on port %d", cfg.APIPort)
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.APIPort),
		Handler:      secureRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
