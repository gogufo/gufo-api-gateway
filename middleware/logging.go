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
