package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"ratelimitergolang/ratelimiterserver"
	"ratelimitergolang/ratelimiterserver/impl"
	"sync"
	"testing"
	"time"
)

// TestRateLimitMiddleware tests the rate limiting middleware
func TestRateLimitMiddleware(t *testing.T) {
	// Create a new SlidingWindowAlgorithm with a small window size for testing
	rateAlgorithm := impl.NewSlidingWindowAlgorithm(100 * time.Millisecond)
	rateAlgorithm.SetEndpointLimit("/api/endpoint1", 2) // Limit 2 requests per 100ms for /api/endpoint1

	// Create a new rate limiter with the rate limiting algorithm
	limiter := impl.NewRateLimiter(rateAlgorithm)

	// Create a new HTTP handler using the rate limiter middleware
	handler := limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	// Send 2 requests to the /api/endpoint1 endpoint within the time window
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest("GET", "/api/endpoint1", nil)
		req.Header.Set("X-Forwarded-For", "127.0.0.1")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if i >= 2 {
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d: Expected status code %d, got %d", i+1, http.StatusTooManyRequests, w.Code)
			} else {
				t.Logf("Request %d: rateLimited with status code %d", i+1, w.Code)
			}
		} else {
			t.Logf("Request %d: OK", i+1)
		}
	}

	// Send another request to the /api/endpoint1 endpoint within the time window, it should be rate limited
	req := httptest.NewRequest("GET", "/api/endpoint1", nil)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Request 3: Expected status code %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}

func TestPerformance(t *testing.T) {
	// Example sliding window rate limiting algorithm with a window size of 1 minute
	rateAlgorithm := impl.NewSlidingWindowAlgorithm(time.Minute)
	rateAlgorithm.SetEndpointLimit("/api/endpoint1", 10) // Limit 10 requests per minute for /api/endpoint1

	// Create rate limiter middleware
	limiter := impl.NewRateLimiter(rateAlgorithm)

	numRequests := 100000
	endpoint := "/api/endpoint1"

	start := time.Now()
	performRequests(numRequests, limiter, endpoint)
	elapsed := time.Since(start)

	fmt.Printf("Performed %d requests in %s\n", numRequests, elapsed)
}

func performRequests(numRequests int, limiter ratelimiterserver.RateLimiter, endpoint string) {
	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", endpoint, nil)
			req.Header.Set("X-Forwarded-For", "127.0.0.1")
			w := httptest.NewRecorder()
			limiter.Limit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, world!"))
			})).ServeHTTP(w, req)
		}()
	}
	wg.Wait()
}
