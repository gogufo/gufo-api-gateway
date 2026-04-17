// Copyright 2019-2026 Alexey Yanchenko <mail@yanchenko.me>
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
	"net/http"
	"strings"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func ProcessPUT(w http.ResponseWriter, r *http.Request, t *pb.Request, version int) {

	// Initialize nested proto structs if missing
	if t.Auth == nil {
		t.Auth = &pb.AuthContext{}
	}
	if t.Context == nil {
		t.Context = &pb.RequestContext{}
	}

	// ===========================
	//  SECURITY CHECK (sign / HMAC / mTLS)
	// ===========================
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {

	case "hmac":
		secret := viper.GetString("security.hmac_secret")
		maxAge := time.Duration(viper.GetInt("security.max_age")) * time.Second

		if t.Auth.Sign == "" || t.Module == "" ||
			!sf.VerifyHMAC(secret, t.Module, t.Auth.Sign, maxAge) {

			errorAnswer(w, r, t, 401, "00001", "Invalid or expired HMAC signature")
			return
		}

	case "sign":
		if t.Auth.Sign == "" || viper.GetString("server.sign") != t.Auth.Sign {
			errorAnswer(w, r, t, 401, "00001", "Invalid signature")
			return
		}

	case "mtls":
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			errorAnswer(w, r, t, 401, "00001", "Client certificate required (mTLS)")
			return
		}

	default:
		errorAnswer(w, r, t, 500, "00002", "Security mode not configured")
		return
	}

	path := r.URL.Path
	patharray := strings.Split(path, "/")
	if len(patharray) < 3 || t.Module == "entrypoint" {
		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Length")
		return
	}

	// 🔐 Check session
	if viper.GetBool("server.session") {
		t = checksession(t, r)
		if t.Auth != nil && t.Auth.Uid != "" && t.Auth.Readonly {
			errorAnswer(w, r, t, 401, "0000235", "Read Only User")
			return
		}
	}

	vrs := "v1"
	t.Context.ApiVersion = vrs
	
	connectgrpc(w, r, t)
}
