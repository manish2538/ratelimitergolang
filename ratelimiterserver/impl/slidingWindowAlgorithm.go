package impl

import (
	"sync"
	"time"
)

// SlidingWindowAlgorithm implements the RateLimitAlgorithm interface using a sliding window approach
type SlidingWindowAlgorithm struct {
	windowSize  time.Duration
	endpointMap map[string]int // Map of endpoint to request limit
	requestMap  map[string]map[string][]time.Time
	mutex       sync.RWMutex
}

// ShouldAllow checks if a request should be allowed based on sliding window rate limiting algorithm
func (sw *SlidingWindowAlgorithm) ShouldAllow(clientIP string, endpoint string) bool {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	limit, ok := sw.endpointMap[endpoint]
	if !ok {
		return true // allow the request if endpoint limit is not configured
	}

	requests, ok := sw.requestMap[endpoint][clientIP]
	if !ok {
		return true // Allow the request if it's the first request from this client to this endpoint
	}

	// Remove requests older than the window size
	for len(requests) > 0 && time.Since(requests[0]) > sw.windowSize {
		requests = requests[1:]
	}
	sw.requestMap[endpoint][clientIP] = requests

	// Check if the number of requests is within the limit
	return len(requests) < limit
}

// RegisterRequest registers a request in the sliding window algorithm
func (sw *SlidingWindowAlgorithm) RegisterRequest(clientIP string, endpoint string) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	if sw.requestMap[endpoint] == nil {
		sw.requestMap[endpoint] = make(map[string][]time.Time)
	}
	sw.requestMap[endpoint][clientIP] = append(sw.requestMap[endpoint][clientIP], time.Now())
}

// SetEndpointLimit sets the rate limit for a specific endpoint
func (sw *SlidingWindowAlgorithm) SetEndpointLimit(endpoint string, limit int) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	sw.endpointMap[endpoint] = limit
}
