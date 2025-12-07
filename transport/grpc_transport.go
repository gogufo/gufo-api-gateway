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
package transport

import (
	"context"
	"fmt"
	"sync"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

// GRPCTransport implements the Transport interface via gRPC calls.
type GRPCTransport struct{}

// cache of host:port for each service
var (
	svcCache sync.Map
)

// Call executes a gRPC call to a remote microservice.
// Host/Port is resolved from cache or via masterservice discovery.
func (t *GRPCTransport) Call(ctx context.Context, svc, method string, req *pb.Request) (*pb.Response, error) {
	host, port := resolveService(svc, req)

	conn, err := sf.GetGRPCConn(
		host, port,
		viper.GetString("security.ca_path"),
		viper.GetString("security.cert_path"),
		viper.GetString("security.key_path"),
		viper.GetBool("security.mtls"),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial failed: %w", err)
	}

	client := pb.NewReverseClient(conn)

	// Timeout per microservice
	timeout := viper.GetDuration(fmt.Sprintf("microservices.%s.timeout", svc))
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := client.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("grpc call failed: %w", err)
	}

	return resp, nil
}

// resolveService returns host and port from cache or asks masterservice.
func resolveService(svc string, req *pb.Request) (string, string) {
	// 1. Check in-memory cache
	if v, ok := svcCache.Load(svc); ok {
		addr := v.(string)
		h, p, _ := splitAddr(addr)
		return h, p
	}

	// 2. Fallback to masterservice discovery (if enabled)
	if viper.GetBool("server.masterservice") && svc != "masterservice" {
		host := viper.GetString("microservices.masterservice.host")
		port := viper.GetString("microservices.masterservice.port")

		ir := &pb.InternalRequest{
			Param:  ptr("getmicroservicebypath"),
			Method: ptr("GET"),
		}
		req.IR = ir
		resp := sf.GRPCConnect(host, port, req)
		if h, ok := resp["host"].(string); ok {
			p := fmt.Sprintf("%v", resp["port"])
			addr := fmt.Sprintf("%s:%s", h, p)
			svcCache.Store(svc, addr)
			return h, p
		}
	}

	// 3. Default to static config
	host := viper.GetString(fmt.Sprintf("microservices.%s.host", svc))
	port := viper.GetString(fmt.Sprintf("microservices.%s.port", svc))
	addr := fmt.Sprintf("%s:%s", host, port)
	svcCache.Store(svc, addr)
	return host, port
}

func splitAddr(addr string) (string, string, error) {
	var host, port string
	n, err := fmt.Sscanf(addr, "%[^:]:%s", &host, &port)
	if n < 2 || err != nil {
		return "", "", fmt.Errorf("invalid addr: %s", addr)
	}
	return host, port, nil
}

func ptr(s string) *string { return &s }
