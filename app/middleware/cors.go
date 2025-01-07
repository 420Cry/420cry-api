package middleware

import (
	"420-api/app/config"
	"log"
	"net/http"
)

var AllowableOrigins []string

// Initialize CORS settings by loading from config
func Initialize() {
	cfg := config.Get()
	if cfg == nil {
		log.Fatal("Config is not initialized")
	}
	AllowableOrigins = []string{
		"https://" + cfg.AllowedOrigin.App,
	}
	log.Printf("Allowable Origins: %v", AllowableOrigins)
}

// CORS Middleware to allow requests from the allowed origins
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("Received Origin: %s", origin)

		// Check if the origin is allowed and set the appropriate CORS headers
		if isAllowedOrigin(origin) {
			log.Printf("Allowed Origin: %s", origin)
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Allow certain methods and headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle OPTIONS method for preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue with the next handler if not a preflight request
		next.ServeHTTP(w, r)
	})
}

// isAllowedOrigin checks if the provided origin is in the AllowableOrigins list.
func isAllowedOrigin(origin string) bool {
	for _, allowedOrigin := range AllowableOrigins {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}
