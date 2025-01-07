package middleware

import (
	"420-api/app/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRateLimiter tests the rate limiting functionality.
func TestRateLimiter(t *testing.T) {
	// Create a new rate limiter that allows 1 request per second with no burst capacity
	rateLimiter := middleware.NewRateLimiter(1, 1)

	// Create a test handler to pass into the rate limiter middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the rate limiter middleware
	middleware := rateLimiter.Middleware(handler)

	// Test that the first request is allowed
	req1, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	middleware.ServeHTTP(rr1, req1)

	// The first request should return a 200 OK status
	if rr1.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr1.Code)
	}

	// Test that the second request is rate-limited (should return 429 Too Many Requests)
	req2, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	middleware.ServeHTTP(rr2, req2)

	// The second request should be rate-limited with a 429 status
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", rr2.Code)
	}

	// Wait for the rate limiter to reset (1 second in this case)
	time.Sleep(1 * time.Second)

	// Test that the third request is allowed after the rate limiter resets
	req3, err := http.NewRequest("GET", "http://localhost/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr3 := httptest.NewRecorder()
	middleware.ServeHTTP(rr3, req3)

	// The third request should return a 200 OK status again
	if rr3.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr3.Code)
	}
}
