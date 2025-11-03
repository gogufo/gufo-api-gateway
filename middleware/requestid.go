package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// RequestID middleware injects a unique X-Request-ID into headers and context.
type RequestID struct{}

// NewRequestID returns a new instance.
func NewRequestID() *RequestID { return &RequestID{} }

type ctxKey string

const requestIDKey ctxKey = "requestID"

// Before adds X-Request-ID if missing.
func (r *RequestID) Before(req *http.Request, ctx context.Context) (context.Context, error) {
	id := req.Header.Get("X-Request-ID")
	if id == "" {
		id = uuid.New().String()
		req.Header.Set("X-Request-ID", id)
	}
	return context.WithValue(ctx, requestIDKey, id), nil
}

// After sets the same header in the response.
func (r *RequestID) After(w http.ResponseWriter, status int, dur time.Duration) {
	if hdr := w.Header().Get("X-Request-ID"); hdr == "" {
		if id, ok := w.(interface{ Header() http.Header }); ok {
			reqID := id.Header().Get("X-Request-ID")
			if reqID != "" {
				w.Header().Set("X-Request-ID", reqID)
			}
		}
	}
}
