// Copyright 2020-2024 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// This file content Handler for API
// Each API function is independend plugin
// and API get reguest in connect with plugin
// Get response from plugin and answer to client
// All data is in JSON format

package handler

import (
	"net/http"
	"time"

	"github.com/gogufo/gufo-api-gateway/middleware"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

var methodHandlers = map[string]func(http.ResponseWriter, *http.Request, *pb.Request, int){
	"OPTIONS": ProcessOPTIONS,
	"GET":     ProcessREQ,
	"HEAD":    ProcessREQ,
	"TRACE":   ProcessREQ,
	"POST":    ProcessREQ,
	"PATCH":   ProcessREQ,
	"DELETE":  ProcessREQ,
	"PUT":     ProcessPUT,
}

// API is the main entrypoint for all REST requests in Gufo Gateway.
// It runs the middleware chain (Before/After) around the appropriate handler.
func API(w http.ResponseWriter, r *http.Request, version int) {
	t := RequestInit(r)

	// 1️⃣ Run global middleware chain (Before)
	ctx, err := middleware.RunBefore(r, r.Context())
	if err != nil {
		errorAnswer(w, r, t, 429, "000429", err.Error())
		return
	}
	start := time.Now()

	// 2️⃣ Core routing by HTTP method
	status := http.StatusOK // default status
	if h, ok := methodHandlers[r.Method]; ok {
		h(w, r.WithContext(ctx), t, version)
	} else {
		ProcessOPTIONS(w, r, t, version)
		status = http.StatusNoContent
	}
	// 3️⃣ Record metrics (QPS + latency)
	ObserveHTTPRequest(r.Method, r.URL.Path, status, start)

	// 4⃣ Run middleware chain (After)
	middleware.RunAfter(w, status, time.Since(start))
}
