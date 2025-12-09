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
