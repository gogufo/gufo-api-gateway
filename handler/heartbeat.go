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
package handler

import (
	"fmt"
	"net/http"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

// Universal heartbeat entry through Gateway.
// Microservices POST here -> Gufo routes to masterservice.
func HeartbeatHandler(w http.ResponseWriter, r *http.Request, t *pb.Request) {

	var payload map[string]interface{}

	ans, err := heartbeatCore(t, payload)
	if err != nil {
		errorAnswer(w, r, t, 500, "0000501", err.Error())
		return
	}

	moduleAnswerv3(w, r, ans, t)
}

// heartbeatCore contains the shared heartbeat business logic.
func heartbeatCore(t *pb.Request, payload map[string]interface{}) (map[string]interface{}, error) {

	msEnabled := viper.GetBool("server.masterservice")

	// ------------------------------------------------------------
	// MODE 2: Standalone mode → return local mock
	// ------------------------------------------------------------
	if !msEnabled {
		mock := map[string]interface{}{
			"leader": true,
			"cron":   true,
			"ttl":    0,
			"epoch":  0,
			"ts":     time.Now().Unix(),
		}
		return mock, nil
	}

	// ------------------------------------------------------------
	// MODE 1: Cluster mode → proxy to MasterService via gRPC
	// ------------------------------------------------------------

	// If payload is nil, initialize empty
	if payload == nil {
		payload = map[string]interface{}{}
	}

	// Always update timestamp before proxying
	payload["ts"] = time.Now().Unix()

	// Convert payload to protobuf.Any
	body, _ := sf.ConvertInterfaceToAny(payload)

	// Build gRPC request for MasterService
	req := &pb.Request{
		Module:  "masterservice",
		Param:   "heartbeat",
		Method:  pb.Method_METHOD_POST,
		Body:    body,
		Context: t.Context,
	}

	host := sf.ConfigString("microservices.masterservice.host")
	port := sf.ConfigString("microservices.masterservice.port")

	// Execute gRPC call to MasterService
	ans := sf.GRPCConnect(host, port, req)

	// Transport-level or application-level failure
	if ans == nil || ans["httpcode"] != nil {
		return nil, fmt.Errorf("masterservice heartbeat failed")
	}

	return ans, nil
}
