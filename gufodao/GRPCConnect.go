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
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
)

// GRPCConnect performs a gRPC call with connection pooling, TLS/mTLS, timeout, and streaming support.
func GRPCConnect(host string, port string, t *pb.Request) map[string]interface{} {
	answer := make(map[string]interface{})

	if host == "" || port == "" {
		answer["httpcode"] = 500
		answer["code"] = "0000238"
		answer["message"] = "Host or Port not specified"
		return answer
	}

	// 🔹 Handle streaming requests
	if t.IR != nil && t.IR.Param != nil && *t.IR.Param == "stream" {
		return GRPCStream(host, port, t)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	// fmt.Fprintln(os.Stderr, ">>> Address:", addr)

	// 🔹 Get connection from pool with TLS/mTLS
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

	// 🔹 Determine timeout per service or default (5s)
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

	answer = ToMapStringInterface(resp.Data)
	copyRequestBack(t, resp.RequestBack)

	return answer
}

// safeModuleName prevents panic if Module is nil
func safeModuleName(t *pb.Request) string {
	if t.Module == nil {
		return "unknown"
	}
	return *t.Module
}

// copyRequestBack updates token/session fields in request
func copyRequestBack(t *pb.Request, rb *pb.Request) {
	if rb == nil {
		return
	}
	if rb.Token != nil {
		t.Token = rb.Token
	}
	if rb.TokenType != nil {
		t.TokenType = rb.TokenType
	}
	if rb.Language != nil {
		t.Language = rb.Language
	}
	if rb.UID != nil {
		t.UID = rb.UID
	}
	if rb.IsAdmin != nil {
		t.IsAdmin = rb.IsAdmin
	}
	if rb.SessionEnd != nil {
		t.SessionEnd = rb.SessionEnd
	}
	if rb.Completed != nil {
		t.Completed = rb.Completed
	}
	if rb.Readonly != nil {
		t.Readonly = rb.Readonly
	}
}

// logOrSentry logs locally or sends to Sentry if enabled
func logOrSentry(err error) {
	if viper.GetBool("server.sentry") {
		sentry.CaptureException(err)
	} else {
		SetErrorLog(err.Error())
	}
}
