package gateway

import (
	"sync"
	"time"
)

// RateLimiter implements a per-IP token bucket rate limiter
type RateLimiter struct {
	limit   int
	buckets map[string]*bucket
	mu      sync.Mutex
}

type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// NewRateLimiter creates a rate limiter with the given requests-per-second limit.
// If limit <= 0, rate limiting is disabled.
func NewRateLimiter(limit int) *RateLimiter {
	rl := &RateLimiter{
		limit:   limit,
		buckets: make(map[string]*bucket),
	}

	// Background cleanup goroutine to prevent memory leaks from stale IPs
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	if rl.limit <= 0 {
		return true // Rate limiting disabled
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[ip]
	if !exists {
		rl.buckets[ip] = &bucket{
			tokens:    float64(rl.limit) - 1,
			lastCheck: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * float64(rl.limit)
	if b.tokens > float64(rl.limit) {
		b.tokens = float64(rl.limit)
	}
	b.lastCheck = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false // Rate limit exceeded — Anti-DoW triggered
}

// cleanup removes stale IP entries every 5 minutes
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for ip, b := range rl.buckets {
			if b.lastCheck.Before(cutoff) {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}
