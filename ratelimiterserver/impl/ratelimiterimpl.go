package impl

import (
	"net/http"
	"ratelimitergolang/ratelimiterserver"
	"time"
)

// rateLimiter implements the RateLimiter interface
type rateLimiter struct {
	rateAlgorithm ratelimiterserver.RateLimitAlgorithm
}

// Limit is the middleware function to enforce rate limits
func (rl *rateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path
		clientIP := r.Header.Get("X-Forwarded-For")

		// Check if the request should be allowed based on the rate limiting algorithm
		if !rl.rateAlgorithm.ShouldAllow(clientIP, endpoint) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Register the request for rate limiting
		rl.rateAlgorithm.RegisterRequest(clientIP, endpoint)

		next.ServeHTTP(w, r)
	})
}

// NewSlidingWindowAlgorithm creates a new SlidingWindowAlgorithm instance
func NewSlidingWindowAlgorithm(windowSize time.Duration) *SlidingWindowAlgorithm {
	return &SlidingWindowAlgorithm{
		windowSize:  windowSize,
		endpointMap: make(map[string]int),
		requestMap:  make(map[string]map[string][]time.Time),
	}
}

func DefaultRateLimiter() *SlidingWindowAlgorithm {
	return NewSlidingWindowAlgorithm(time.Minute)
}

// NewRateLimiter creates a new rate limiter with the given rate limiting algorithm
func NewRateLimiter(rateAlgorithm ratelimiterserver.RateLimitAlgorithm) ratelimiterserver.RateLimiter {
	return &rateLimiter{
		rateAlgorithm: rateAlgorithm,
	}
}
