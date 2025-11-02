// Copyright 2024-2025 Alexey Yanchenko <mail@yanchenko.me>
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

package gufodao

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type connItem struct {
	conn   *grpc.ClientConn
	expiry time.Time
}

var (
	connPool sync.Map
	ttl      = 5 * time.Minute
)

// GetGRPCConn returns a cached gRPC connection (TLS or mTLS) with keepalive and timeout.
// If no active connection exists, a new one is established and stored in the pool.
func GetGRPCConn(host, port, ca, cert, key string, useMTLS bool) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	// 🔹 Check existing connection in pool
	if v, ok := connPool.Load(addr); ok {
		item := v.(connItem)
		if time.Now().Before(item.expiry) {
			return item.conn, nil
		}
		_ = item.conn.Close()
		connPool.Delete(addr)
	}

	// 🔹 Setup TLS credentials
	var creds credentials.TransportCredentials
	if useMTLS {
		certificate, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %w", err)
		}

		caCert, err := os.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA file: %w", err)
		}

		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certs")
		}

		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      caPool,
		})
	} else {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}

	// 🔹 Define keepalive parameters for long-lived connections
	kaParams := keepalive.ClientParameters{
		Time:                30 * time.Second, // send pings every 30s
		Timeout:             10 * time.Second, // wait 10s for ack
		PermitWithoutStream: true,             // allow ping when idle
	}

	// 🔹 Establish new connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(kaParams),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC %s: %w", addr, err)
	}

	// 🔹 Store connection in pool
	connPool.Store(addr, connItem{
		conn:   conn,
		expiry: time.Now().Add(ttl),
	})

	return conn, nil
}
