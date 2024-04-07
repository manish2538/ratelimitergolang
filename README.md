# Rate Limiter
A simple rate limiter implementation in Go for HTTP servers.

# Overview
This project provides a middleware package for rate limiting HTTP requests based on client IP address and endpoint. It uses a sliding window algorithm to enforce rate limits for specific endpoints.

# Installation
To use this package in your Go project, simply import it:

```
import "github.com/manish2538/ratelimiter/impl"
```

Then, run go get to install the package:

```
go get -u github.com/manish2538/ratelimiter/impl
```

# Usage
Here's an example of how to use the rate limiter middleware in your HTTP server:

```
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/manish2538/ratelimiter/impl"
)

const (
	DEFAULT_PORT = ":8080"
)

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
	fmt.Printf("Starting server at %s port\n", DEFAULT_PORT)
	http.ListenAndServe(DEFAULT_PORT, limiter.Limit(handler))
}

func initRateLimiter() impl.RateLimiter {
	rateAlgorithm := impl.DefaultRateLimiter()
	// Set rate limits for specific endpoints
	rateAlgorithm.SetEndpointLimit("/api/user", 3)    // Limit 3 requests per minute for /api/endpoint1
	rateAlgorithm.SetEndpointLimit("/api/profile", 4) // Limit 4 requests per minute for /api/endpoint2

	return impl.NewRateLimiter(rateAlgorithm)
}

```
# Configuration
You can configure rate limits for specific endpoints by calling the SetEndpointLimit method of the rate limiter algorithm. The example above demonstrates how to set rate limits for two endpoints (/api/user and /api/profile).

# Contributing
Contributions are welcome! Please feel free to submit bug reports, feature requests, or pull requests.

# License
This project is licensed under the MIT License - see the LICENSE file for details.

Feel free to modify the content as needed to fit your project's specific requirements and conventions.
