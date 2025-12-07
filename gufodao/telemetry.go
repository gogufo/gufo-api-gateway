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
package gufodao

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var tracerProvider *sdktrace.TracerProvider

// InitTelemetry initializes OpenTelemetry tracer with optional OTLP exporter.
func InitTelemetry(ctx context.Context) {
	endpoint := os.Getenv("GUFO_OTEL_ENDPOINT")
	service := os.Getenv("GUFO_OTEL_SERVICE_NAME")
	if service == "" {
		service = "gufo-gateway"
	}

	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
		),
	)

	if endpoint == "" {
		// No OTLP endpoint -> use no-op tracer
		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
		)
		otel.SetTracerProvider(tracerProvider)
		log.Println("[telemetry] running in local mode (no OTLP export)")
		return
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // set false if using HTTPS
	)
	if err != nil {
		log.Printf("[telemetry] failed to init OTLP exporter: %v", err)
		return
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	log.Printf("[telemetry] OpenTelemetry connected to %s", endpoint)
}

// ShutdownTelemetry flushes and closes the tracer provider.
func ShutdownTelemetry(ctx context.Context) {
	if tracerProvider != nil {
		_ = tracerProvider.Shutdown(ctx)
	}
}
