package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Logger middleware records request method, path, status, and latency.
type Logger struct{}

// NewLogger returns a new structured logger middleware.
func NewLogger() *Logger { return &Logger{} }

func (l *Logger) Before(r *http.Request, ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, "startTime", time.Now())
	return ctx, nil
}

func (l *Logger) After(w http.ResponseWriter, status int, dur time.Duration) {
	log.Printf("[REQ] %s status=%d duration=%v",
		w.Header().Get("X-Request-ID"),
		status,
		dur,
	)

}
