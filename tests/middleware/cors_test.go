package middleware_test

import (
	"420-api/app/config"
	"420-api/app/middleware"
	"420-api/app/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup test configuration for CORS
func setupCORS() {
	cfg := &types.Config{
		APIPort: 8080,
		AllowedOrigin: types.AllowedOrigin{
			App: "example.com",
		},
	}
	config.Set(cfg)
	middleware.Initialize()
}

func TestCORS_AllowOrigin(t *testing.T) {
	setupCORS()

	// Create a test handler
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a new request with the allowed origin header
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://example.com")

	// Record the response
	rr := httptest.NewRecorder()
	handler = middleware.CORS(handler)
	handler.ServeHTTP(rr, req)

	// Check the response status and headers
	assert.Equal(t, http.StatusOK, rr.Code)
	allowedOrigin := rr.Header().Get("Access-Control-Allow-Origin")
	assert.Equal(t, "https://example.com", allowedOrigin)
}

func TestCORS_DisallowOrigin(t *testing.T) {
	setupCORS()

	// Create a test handler
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create a new request with a disallowed origin header
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://unauthorized.com")

	// Record the response
	rr := httptest.NewRecorder()
	handler = middleware.CORS(handler)
	handler.ServeHTTP(rr, req)

	// Check the response status and headers
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
}
