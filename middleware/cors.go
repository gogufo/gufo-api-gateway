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
