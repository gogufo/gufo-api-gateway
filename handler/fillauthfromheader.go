// Copyright 2026 Alexey Yanchenko <mail@yanchenko.me>
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
	"net/http"
	"strings"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func fillAuthFromHeaders(t *pb.Request, r *http.Request) *pb.Request {
	if t.Auth == nil {
		t.Auth = &pb.AuthContext{}
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader != "" {
		t.Auth.Token = authHeader

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 {
			t.Auth.TokenType = strings.TrimSpace(parts[0])
			t.Auth.Token = strings.TrimSpace(parts[1])
		}
	}

	apiToken := strings.TrimSpace(r.Header.Get("X-API-Token"))
	if apiToken == "" {
		apiToken = strings.TrimSpace(r.Header.Get("X-Api-Token"))
	}
	if apiToken == "" {
		apiToken = strings.TrimSpace(r.Header.Get("X-API-TOKEN"))
	}

	if apiToken != "" {
		// если в proto нет отдельного поля ApiToken — кладём в Token
		// но только если Authorization не был передан
		if t.Auth.Token == "" {
			t.Auth.Token = apiToken
			t.Auth.TokenType = "ApiToken"
		}

		// если есть Meta в RequestContext — сохраним ещё и там
		if t.Context == nil {
			t.Context = &pb.RequestContext{}
		}
		if t.Context.Meta == nil {
			t.Context.Meta = map[string]string{}
		}
		t.Context.Meta["x_api_token"] = apiToken
	}

	return t
}
