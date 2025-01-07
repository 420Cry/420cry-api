package middleware

import (
	"420-api/app/config"
	"net/http"
)

var AllowableOrigins []string

func Initialize() {
	cfg := config.Load()
	AllowableOrigins = []string{
		"https://" + cfg.ALLOWED_ORIGIN.APP,
	}
}

// CORS Middleware to allow requests from the allowed origins
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range AllowableOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}

		// Allow certain methods and headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If the method is OPTIONS, return early
		if r.Method == http.MethodOptions {
			return
		}

		// Continue with the next handler
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

// handleCORSPreflight responds to preflight OPTIONS requests for CORS
func handleCORSPreflight(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if !isAllowedOrigin(origin) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Set CORS headers for the OPTIONS preflight request
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}
