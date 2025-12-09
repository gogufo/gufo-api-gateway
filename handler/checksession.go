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
// This file content Handler for API
// Each API function is independend plugin
// and API get reguest in connect with plugin
// Get response from plugin and answer to client
// All data is in JSON format

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

// checksession validates the session of an incoming request.
// It extracts the token from the Authorization header (Bearer format preferred)
// or, for backward compatibility, from query parameters (?access_token=...).
// Then it verifies the session against the Session microservice (via Masterservice or direct host).
func checksession(t *pb.Request, r *http.Request) *pb.Request {
	p := bluemonday.UGCPolicy()
	var tokenHeader string

	// 1) Extract token from Authorization header (RFC 6750)
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenHeader = parts[1]
		} else {
			sf.SetErrorLog("checksession: invalid Authorization header format")
			return t
		}
	} else {
		// 2) Legacy fallback: access_token in URL query
		if q := r.URL.Query().Get("access_token"); q != "" {
			tokenHeader = p.Sanitize(q)
			if tt := r.URL.Query().Get("token_type"); tt != "" {
				tokenHeader = fmt.Sprintf("%s %s", p.Sanitize(tt), tokenHeader)
			} else {
				tokenHeader = "Bearer " + tokenHeader
			}
		}
	}

	// 3) No token found — return without session data
	if tokenHeader == "" {
		return t
	}

	// 4) Determine session microservice host
	var host, port string
	if viper.GetBool("server.masterservice") {
		host = viper.GetString("microservices.masterservice.host")
		port = viper.GetString("microservices.masterservice.port")

		mst := &pb.InternalRequest{
			Param:  sf.StringPtr("getsessionhost"),
			Method: sf.StringPtr("GET"),
		}
		t.IR = mst
		t.Token = &tokenHeader

		ans := sf.GRPCConnect(host, port, t)
		if ans["httpcode"] != nil {
			return t // masterservice error
		}

		host = fmt.Sprintf("%v", ans["host"])
		port = fmt.Sprintf("%v", ans["port"])
	} else {
		if !viper.IsSet("microservices.session.host") {
			return t
		}
		host = viper.GetString("microservices.session.host")
		port = viper.GetString("microservices.session.port")
	}

	// 5) Call Session microservice to validate the token
	mstb := &pb.InternalRequest{
		Param:  sf.StringPtr("checksession"),
		Method: sf.StringPtr("GET"),
	}
	t.IR = mstb
	t.Token = &tokenHeader

	ans := sf.GRPCConnect(host, port, t)
	if ans["error"] != nil {
		sf.SetErrorLog(fmt.Sprintf("checksession: gRPC error: %v", ans["error"]))
		return t
	}

	// 6) Populate response fields from Session service
	if v := ans["uid"]; v != nil {
		uid := fmt.Sprintf("%v", v)
		t.UID = &uid
	}
	if v := ans["isadmin"]; v != nil {
		i, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		t.IsAdmin = sf.Int32Ptr(int32(i))
	}
	if v := ans["sessionend"]; v != nil {
		i, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		t.SessionEnd = sf.Int32Ptr(int32(i))
	}
	if v := ans["completed"]; v != nil {
		i, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		t.Completed = sf.Int32Ptr(int32(i))
	}
	if v := ans["readonly"]; v != nil {
		i, _ := strconv.Atoi(fmt.Sprintf("%v", v))
		t.Readonly = sf.Int32Ptr(int32(i))
	}
	if v := ans["token"]; v != nil {
		tkn := fmt.Sprintf("%v", v)
		t.Token = &tkn
	}
	if v := ans["token_type"]; v != nil {
		tkntp := fmt.Sprintf("%v", v)
		t.TokenType = &tkntp
	}

	return t
}
