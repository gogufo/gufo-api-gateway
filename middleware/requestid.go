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
