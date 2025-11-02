// Copyright 2025 Alexey Yanchenko
// Part of Gufo API Gateway
//
// Transport layer abstraction for calling microservices via different backends.
// Current default: gRPC Transport (PR-5)

package transport

import (
	"context"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

// Transport defines a unified interface for all transports (gRPC, HTTP, MQ, etc.)
type Transport interface {
	Call(ctx context.Context, svc string, method string, req *pb.Request) (*pb.Response, error)
}

// DefaultTransport holds the active transport instance (initialized in main)
var DefaultTransport Transport

// Register sets the current transport implementation.
func Register(t Transport) {
	DefaultTransport = t
}

// Get returns the registered transport (default: gRPCTransport).
func Get() Transport {
	return DefaultTransport
}
