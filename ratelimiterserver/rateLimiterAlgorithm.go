package ratelimiterserver

// RateLimitAlgorithm defines the interface for rate limiting algorithms
type RateLimitAlgorithm interface {
	ShouldAllow(clientIP string, endpoint string) bool
	RegisterRequest(clientIP string, endpoint string)
	SetEndpointLimit(endpoint string, limit int)
}
