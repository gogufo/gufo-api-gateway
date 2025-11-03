package middleware

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// RateLimiter provides simple in-memory token bucket limiting.
type RateLimiter struct {
	capacity int
	tokens   int
	refill   time.Duration
	mu       sync.Mutex
	last     time.Time
}

// NewRateLimiter creates a new limiter (e.g., 100 req/s).
func NewRateLimiter(capacity int, refill time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity: capacity,
		tokens:   capacity,
		refill:   refill,
		last:     time.Now(),
	}
}

func (rl *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.last)
	tokensToAdd := int(elapsed / rl.refill)
	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.last = now
	}
}

func (rl *RateLimiter) Before(r *http.Request, ctx context.Context) (context.Context, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refillTokens()
	if rl.tokens <= 0 {
		return ctx, errors.New("rate limit exceeded")
	}
	rl.tokens--
	return ctx, nil
}

func (rl *RateLimiter) After(w http.ResponseWriter, status int, dur time.Duration) {}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
