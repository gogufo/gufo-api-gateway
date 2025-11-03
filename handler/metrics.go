// Copyright 2025 Alexey Yanchenko <mail@yanchenko.me>
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

package handler

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// -------------------------
// Metric definitions
// -------------------------

var (
	// Total number of HTTP requests (QPS counter)
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gufo_http_requests_total",
			Help: "Total number of HTTP requests received, labeled by method and path.",
		},
		[]string{"method", "path"},
	)

	// Request duration histogram (for latency analysis)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gufo_http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations (seconds).",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// gRPC connection pool metrics (optional for future use)
	grpcPoolHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gufo_grpc_pool_hits_total",
			Help: "Number of successful gRPC connection pool hits.",
		},
	)
	grpcPoolMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gufo_grpc_pool_misses_total",
			Help: "Number of gRPC connection pool misses.",
		},
	)
)

// -------------------------
// Initialization
// -------------------------

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(grpcPoolHits)
	prometheus.MustRegister(grpcPoolMisses)
}

// -------------------------
// Public API
// -------------------------

// ObserveHTTPRequest records request metrics.
func ObserveHTTPRequest(method, path string, status int, start time.Time) {
	httpRequestsTotal.WithLabelValues(method, path).Inc()
	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues(method, path, http.StatusText(status)).Observe(duration)
}

// IncrementGRPCPoolHit increments the gRPC pool hit counter.
func IncrementGRPCPoolHit() {
	grpcPoolHits.Inc()
}

// IncrementGRPCPoolMiss increments the gRPC pool miss counter.
func IncrementGRPCPoolMiss() {
	grpcPoolMisses.Inc()
}

// MetricsHandler exposes all registered Prometheus metrics.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
