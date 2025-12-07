// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	rps    int
	tokens int
	burst  int
	refill time.Duration
	mu     sync.Mutex
	last   time.Time
}

// NewRateLimiter creates a new limiter (e.g., 100 req/s).
func NewRateLimiter(rps int, refill time.Duration, burst int) *RateLimiter {
	if rps <= 0 {
		rps = 100
	}
	if burst <= rps {
		burst = rps * 2
	}

	return &RateLimiter{
		rps:    rps,
		tokens: burst, // start with full bucket
		burst:  burst,
		refill: refill,
		last:   time.Now(),
	}
}

func (rl *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.last)
	tokensToAdd := int(elapsed / rl.refill)

	if tokensToAdd > 0 {
		rl.tokens += tokensToAdd
		if rl.tokens > rl.burst {
			rl.tokens = rl.burst
		}
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
