// Copyright 2025 Alexey Yanchenko <mail@yanchenko.me>
//
// Part of the Gufo API Gateway project.
// Licensed under the Apache License, Version 2.0.

package middleware

import (
	"context"
	"net/http"
	"time"
)

// Middleware defines the interface for request lifecycle hooks.
type Middleware interface {
	// Before runs before the request is sent to the microservice.
	// It may enrich the context or reject the request with an error.
	Before(r *http.Request, ctx context.Context) (context.Context, error)

	// After runs after the request completes.
	// It may log results, collect metrics, or modify headers.
	After(w http.ResponseWriter, status int, dur time.Duration)
}

var chain []Middleware

// Register adds a middleware to the global execution chain.
func Register(m Middleware) {
	chain = append(chain, m)
}

// RunBefore executes all registered middleware Before() in order.
func RunBefore(r *http.Request, ctx context.Context) (context.Context, error) {
	var err error
	for _, m := range chain {
		ctx, err = m.Before(r, ctx)
		if err != nil {
			return ctx, err
		}
	}
	return ctx, nil
}

// RunAfter executes all registered middleware After() in reverse order.
func RunAfter(w http.ResponseWriter, status int, dur time.Duration) {
	for i := len(chain) - 1; i >= 0; i-- {
		chain[i].After(w, status, dur)
	}
}
