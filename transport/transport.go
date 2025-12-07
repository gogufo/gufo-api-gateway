// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
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
// Transport layer abstraction for calling microservices via different backends.

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
