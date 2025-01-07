package main

import (
	"420-api/app/config"
	"420-api/app/middleware"
	"420-api/app/routes"
	"log"
	"net/http"
	"strconv"

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
	log.Printf("Server started on port %d", cfg.API_PORT)
	err := http.ListenAndServe(":"+strconv.Itoa(cfg.API_PORT), secureRouter)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
