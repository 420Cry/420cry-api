package middleware_test

import (
	"420-api/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Setup a basic handler to test middleware
func basicHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestSecurityHeaders_AllowOrigin(t *testing.T) {
	setupCORS()
	// Create a test handler to pass into the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a request with a mocked Origin header
	req, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Origin", "https://example.com") // Mock the Origin header here

	// Record the response
	rr := httptest.NewRecorder()
	// Wrap the handler with the SecurityHeaders middleware
	middleware := middleware.SecurityHeaders(handler)

	// Serve the HTTP request with the middleware
	middleware.ServeHTTP(rr, req)

	// Check if the response code is 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Check if the Access-Control-Allow-Origin header is set correctly
	expectedHeader := "https://example.com"
	if rr.Header().Get("Access-Control-Allow-Origin") != expectedHeader {
		t.Errorf("expected Access-Control-Allow-Origin to be %s, got %s", expectedHeader, rr.Header().Get("Access-Control-Allow-Origin"))
	}

	// Check other headers
	expectedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	if rr.Header().Get("Access-Control-Allow-Methods") != expectedMethods {
		t.Errorf("expected Access-Control-Allow-Methods to be %s, got %s", expectedMethods, rr.Header().Get("Access-Control-Allow-Methods"))
	}

	expectedHeaders := "Content-Type, Authorization"
	if rr.Header().Get("Access-Control-Allow-Headers") != expectedHeaders {
		t.Errorf("expected Access-Control-Allow-Headers to be %s, got %s", expectedHeaders, rr.Header().Get("Access-Control-Allow-Headers"))
	}

	// Check other security headers
	expectedCSP := "default-src 'self'"
	if rr.Header().Get("Content-Security-Policy") != expectedCSP {
		t.Errorf("expected Content-Security-Policy to be %s, got %s", expectedCSP, rr.Header().Get("Content-Security-Policy"))
	}

	expectedXContentTypeOptions := "nosniff"
	if rr.Header().Get("X-Content-Type-Options") != expectedXContentTypeOptions {
		t.Errorf("expected X-Content-Type-Options to be %s, got %s", expectedXContentTypeOptions, rr.Header().Get("X-Content-Type-Options"))
	}

	expectedXFrameOptions := "DENY"
	if rr.Header().Get("X-Frame-Options") != expectedXFrameOptions {
		t.Errorf("expected X-Frame-Options to be %s, got %s", expectedXFrameOptions, rr.Header().Get("X-Frame-Options"))
	}

	expectedStrictTransportSecurity := "max-age=63072000; includeSubDomains"
	if rr.Header().Get("Strict-Transport-Security") != expectedStrictTransportSecurity {
		t.Errorf("expected Strict-Transport-Security to be %s, got %s", expectedStrictTransportSecurity, rr.Header().Get("Strict-Transport-Security"))
	}
}

func TestSecurityHeaders_DisallowOrigin(t *testing.T) {
	// Test disallowed origin
	handler := http.Handler(http.HandlerFunc(basicHandler))
	// Initialize middleware
	handler = middleware.SecurityHeaders(handler)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://unauthorized.com")

	// Record the response
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert Forbidden response
	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestSecurityHeaders_PreflightRequest(t *testing.T) {
	setupCORS()
	// Test CORS preflight request (OPTIONS method)
	handler := http.Handler(http.HandlerFunc(basicHandler))
	// Initialize middleware
	handler = middleware.SecurityHeaders(handler)

	req, err := http.NewRequest(http.MethodOptions, "/", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "https://example.com")

	// Record the response
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert status and CORS headers
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
}

func TestSecurityHeaders_MissingOrigin(t *testing.T) {
	// Test missing origin header
	handler := http.Handler(http.HandlerFunc(basicHandler))
	// Initialize middleware
	handler = middleware.SecurityHeaders(handler)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	// Record the response
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Assert Forbidden response (since Origin is missing)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}
