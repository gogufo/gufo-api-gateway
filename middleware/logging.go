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
