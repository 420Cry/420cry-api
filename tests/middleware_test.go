package tests

import (
	"420-api/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	// Initialize CORS configuration
	middleware.Initialize() // Ensure that AllowableOrigins is correctly set

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with CORS middleware
	r := mux.NewRouter()
	r.Use(middleware.CORS)
	r.Handle("/", handler)

	// Test allowed origin
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://allowed-origin.com")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "https://allowed-origin.com", rr.Header().Get("Access-Control-Allow-Origin"))

	// Test disallowed origin
	req.Header.Set("Origin", "https://disallowed-origin.com")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSecurityHeaders(t *testing.T) {
	// Initialize CORS configuration
	middleware.Initialize() // Ensure that the headers are properly set

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with SecurityHeaders middleware
	r := mux.NewRouter()
	r.Use(middleware.SecurityHeaders)
	r.Handle("/", handler)

	// Test request with allowed origin
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://allowed-origin.com")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check security headers
	assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rr.Header().Get("X-Frame-Options"))
	assert.Equal(t, "max-age=63072000; includeSubDomains", rr.Header().Get("Strict-Transport-Security"))
}

func TestRateLimiting(t *testing.T) {
	// Create a rate limiter with 1 request per 5 seconds
	limiter := middleware.NewRateLimiter(1, 5)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a test server with Rate Limiting middleware
	r := mux.NewRouter()
	r.Use(limiter.Middleware)
	r.Handle("/", handler)

	// Simulate a request within the limit
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Simulate another request within the limit
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Simulate a request exceeding the limit
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
}
