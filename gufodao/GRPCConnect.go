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
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
)

// GRPCConnect performs a gRPC call with connection pooling, TLS/mTLS and timeout.
func GRPCConnect(host string, port string, t *pb.Request) map[string]interface{} {
	answer := make(map[string]interface{})

	if host == "" || port == "" {
		answer["httpcode"] = 500
		answer["code"] = "0000238"
		answer["message"] = "Host or Port not specified"
		return answer
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	// 🔹 Get connection from pool
	conn, err := GetGRPCConn(
		host,
		port,
		viper.GetString("security.ca_path"),
		viper.GetString("security.cert_path"),
		viper.GetString("security.key_path"),
	)
	if err != nil {
		logOrSentry(fmt.Errorf("grpc dial failed for %s: %w", addr, err))
		answer["httpcode"] = 400
		answer["code"] = "0000234"
		answer["message"] = err.Error()
		return answer
	}

	client := pb.NewReverseClient(conn)

	// 🔹 Timeout per service
	timeout := viper.GetDuration(fmt.Sprintf("microservices.%s.timeout", safeModuleName(t)))
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 🔹 Perform RPC
	resp, err := client.Do(ctx, t)
	if err != nil {
		logOrSentry(fmt.Errorf("grpc call failed for %s: %w", addr, err))
		answer["httpcode"] = 500
		answer["code"] = "0000236"
		answer["message"] = fmt.Sprintf("Module connection error: %s", err.Error())
		return answer
	}

	// 🔹 Convert Data (Any → interface{})
	if resp.Data != nil {
		answer = ToMapStringInterface(resp.Data)
	}

	// 🔹 Propagate microservice error
	if resp.Error != nil {
		answer["error"] = map[string]interface{}{
			"code":      resp.Error.Code,
			"key":       resp.Error.Key,
			"message":   resp.Error.Message,
			"retryable": resp.Error.Retryable,
		}
	}

	// 🔹 Attach meta (optional)
	if resp.Meta != nil {
		answer["meta"] = map[string]interface{}{
			"status":     resp.Meta.Status,
			"trace_id":   resp.Meta.TraceId,
			"request_id": resp.Meta.RequestId,
			"node":       resp.Meta.Node,
		}
	}

	// 🔹 Sync auth/context back into request
	copyRequestBack(t, resp.RequestBack)

	return answer
}

// safeModuleName prevents panic
func safeModuleName(t *pb.Request) string {
	if t == nil || t.Module == "" {
		return "unknown"
	}
	return t.Module
}

// copyRequestBack updates Auth and Context from response
func copyRequestBack(t *pb.Request, rb *pb.Request) {
	if rb == nil {
		return
	}

	if rb.Auth != nil {
		t.Auth = rb.Auth
	}

	if rb.Context != nil {
		t.Context = rb.Context
	}
}

// logOrSentry logs locally or sends to Sentry
func logOrSentry(err error) {
	if viper.GetBool("server.sentry") {
		sentry.CaptureException(err)
	} else {
		SetErrorLog(err.Error())
	}
}
