package utils

import (
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiterMap = make(map[string]*rate.Limiter)
	mu         sync.Mutex
)

// GetRateLimiter returns a rate limiter for a given user ID
func GetRateLimiter(userID string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := limiterMap[userID]; exists {
		return limiter
	}

	// Create a new rate limiter for the user if it doesn't exist
	limiter := rate.NewLimiter(5, 5)
	limiterMap[userID] = limiter
	return limiter
}
