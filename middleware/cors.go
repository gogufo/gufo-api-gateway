package middleware

import (
	"context"
	"net/http"
	"time"
)

// CORS middleware adds Access-Control headers.
type CORS struct{}

func NewCORS() *CORS { return &CORS{} }

func (c *CORS) Before(r *http.Request, ctx context.Context) (context.Context, error) {
	if r.Method == http.MethodOptions {
		return ctx, nil
	}
	return ctx, nil
}

func (c *CORS) After(w http.ResponseWriter, status int, dur time.Duration) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
}
