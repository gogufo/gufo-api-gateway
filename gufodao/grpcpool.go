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

package gufodao

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// connItem stores an active gRPC connection and its expiration timestamp.
type connItem struct {
	conn   *grpc.ClientConn
	expiry time.Time
}

var (
	connPool sync.Map          // connection pool: addr -> connItem
	ttl      = 5 * time.Minute // connection TTL before re-dial
)

// init launches a background sweeper that periodically removes expired connections.
func init() {
	go poolSweeper()
}

// 🧹 poolSweeper runs every 10 minutes to clean expired connections.
func poolSweeper() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		connPool.Range(func(key, value any) bool {
			item := value.(connItem)
			if now.After(item.expiry) {
				item.conn.Close()
				connPool.Delete(key)
			}
			return true
		})
	}
}

// GetGRPCConn returns a pooled gRPC connection with TLS/mTLS, keepalive, and retry policy.
func GetGRPCConn(host, port, ca, cert, key string) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	// 1️⃣ Check existing connection
	if v, ok := connPool.Load(addr); ok {
		item := v.(connItem)
		if time.Now().Before(item.expiry) {
			return item.conn, nil
		}
		item.conn.Close()
		connPool.Delete(addr)
	}

	// 2️⃣ Prepare transport credentials
	var creds credentials.TransportCredentials
	mode := strings.ToLower(viper.GetString("security.mode"))

	if mode == "mtls" {
		certificate, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("load keypair: %w", err)
		}
		caCert, err := os.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("read CA file: %w", err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      caPool,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	// 3️⃣ Keepalive parameters
	kaParams := keepalive.ClientParameters{
		Time:                30 * time.Second, // ping every 30s
		Timeout:             10 * time.Second, // wait 10s for ack
		PermitWithoutStream: true,
	}

	// 4️⃣ Retry policy
	retrySC := `{
		"methodConfig": [{
			"name": [{"service": "Reverse"}],
			"retryPolicy": {
				"MaxAttempts": 4,
				"InitialBackoff": "0.2s",
				"MaxBackoff": "2s",
				"BackoffMultiplier": 1.6,
				"RetryableStatusCodes": ["UNAVAILABLE","RESOURCE_EXHAUSTED"]
			}
		}]
	}`

	// 5️⃣ Dial with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Fprintln(os.Stderr, ">>> GRPC DIAL ADDR =", addr)
	fmt.Fprintln(os.Stderr, ">>> GRPC DIAL OPTS =", grpc.WithTransportCredentials(creds))

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(kaParams),
		grpc.WithDefaultServiceConfig(retrySC),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, ">>> GRPC DIAL ERROR =", err.Error())
		return nil, fmt.Errorf("dial %s failed: %w", addr, err)
	}

	// 6️⃣ Store in pool
	connPool.Store(addr, connItem{
		conn:   conn,
		expiry: time.Now().Add(ttl),
	})

	fmt.Fprintln(os.Stderr, ">>> GRPC DIAL OK =", addr)

	return conn, nil
}
