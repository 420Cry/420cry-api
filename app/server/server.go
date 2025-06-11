// server main. this should be refactored
package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"cry-api/app/api"
	"cry-api/app/config"
	"cry-api/app/database"

	"github.com/gorilla/mux"
)

func main() {
	// Load the configuration settings
	cfg := config.Get()
	origin := cfg.CryAppURL
	dbConn, err := database.GetDBConnection()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	// Set up the router
	r := mux.NewRouter()
	db := dbConn.GetDB()

	// Register all routes dynamically
	api.RegisterAllRoutes(r, db)

	// Wrap router with CORS middleware
	corsRouter := enableCORS(r, origin)

	// Define the HTTP server
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.APIPort),
		Handler:           corsRouter,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start the server and check for errors
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// CORS middleware to allow only a specific origin
func enableCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "", http.StatusForbidden)
	})
}
