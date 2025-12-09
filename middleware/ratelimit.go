// Copyright 2019-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Business Source License 1.1 (the "License");
// you may not use this file except in compliance with the License.
//
// You may obtain a copy of the License in the LICENSE file at the root of this repository.
//
// As of the Change Date specified in that file, in accordance with the Business Source
// License, use of this software will be governed by the Apache License, Version 2.0.
//
// THIS SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NON-INFRINGEMENT.
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
