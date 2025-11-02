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

// Copyright 2024-2025 Alexey Yanchenko <mail@yanchenko.me>

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
	kaParams = keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}
	retrySC = `{
	  "methodConfig": [{
		"name": [{"service": "proto.Reverse"}],
		"retryPolicy": {
		  "MaxAttempts": 4,
		  "InitialBackoff": "0.2s",
		  "MaxBackoff": "2s",
		  "BackoffMultiplier": 2.0,
		  "RetryableStatusCodes": ["UNAVAILABLE","DEADLINE_EXCEEDED"]
		}
	  }]
	}`
)

func init() {
	// background GC for expired connections
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			connPool.Range(func(k, v any) bool {
				item := v.(connItem)
				if time.Now().After(item.expiry) {
					item.conn.Close()
					connPool.Delete(k)
				}
				return true
			})
		}
	}()
}

// GetGRPCConn returns cached gRPC connection (TLS or mTLS) with keepalive and retries.
func GetGRPCConn(host, port, ca, cert, key string, useMTLS bool) (*grpc.ClientConn, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	if v, ok := connPool.Load(addr); ok {
		item := v.(connItem)
		if time.Now().Before(item.expiry) {
			return item.conn, nil
		}
		item.conn.Close()
		connPool.Delete(addr)
	}

	var creds credentials.TransportCredentials
	if useMTLS {
		certificate, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		caCert, err := os.ReadFile(ca)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      caPool,
		})
	} else {
		creds = credentials.NewClientTLSFromCert(nil, "")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(kaParams),
		grpc.WithDefaultServiceConfig(retrySC),
	)
	if err != nil {
		return nil, err
	}

	connPool.Store(addr, connItem{conn: conn, expiry: time.Now().Add(ttl)})
	return conn, nil
}
