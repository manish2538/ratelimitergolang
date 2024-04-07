package main

import (
	"fmt"
	"net/http"
	"ratelimitergolang/ratelimiterserver"
	"ratelimitergolang/ratelimiterserver/impl"
)

const (
	DEFAULT_PORT = ":8080"
)

func initRateLimiter() ratelimiterserver.RateLimiter {
	rateAlgorithm := impl.DefaultRateLimiter()
	rateAlgorithm.SetEndpointLimit("/api/user", 3)    // Limit 3 requests per minute for /api/endpoint1
	rateAlgorithm.SetEndpointLimit("/api/profile", 4) // Limit 4 requests per minute for /api/endpoint2
	return impl.NewRateLimiter(rateAlgorithm)
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	})
	// Create rate limiter middleware
	limiter := initRateLimiter()
	// Create HTTP server with rate limiting middleware
	fmt.Printf("Starting server at %s port", DEFAULT_PORT)
	http.ListenAndServe(DEFAULT_PORT, limiter.Limit(handler))
}
