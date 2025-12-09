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
	"encoding/json"
	"net/http"
	"strings"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/spf13/viper"
)

func ProcessREQ(w http.ResponseWriter, r *http.Request, t *pb.Request, version int) {

	// ===========================
	//  SECURITY CHECK (same as gRPC Do)
	// ===========================
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {
	case "hmac":
		secret := viper.GetString("security.hmac_secret")
		maxAge := time.Duration(viper.GetInt("security.max_age")) * time.Second

		if t.Sign == nil || t.Module == nil ||
			!sf.VerifyHMAC(secret, *t.Module, *t.Sign, maxAge) {

			errorAnswer(w, r, t, 401, "00001", "Invalid or expired HMAC signature")
			return
		}

	case "sign":
		if t.Sign == nil || viper.GetString("server.sign") != *t.Sign {
			errorAnswer(w, r, t, 401, "00001", "Invalid signature")
			return
		}

	case "mtls":
		// r.TLS != nil only when HTTPS with client certificate
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			errorAnswer(w, r, t, 401, "00001", "Client certificate required (mTLS)")
			return
		}

	default:
		errorAnswer(w, r, t, 500, "00002", "Security mode not configured")
		return
	}

	//Determinate plugin name, params etc.
	//
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	//p := bluemonday.UGCPolicy()

	if pathlenth < 3 {

		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Lenth")

		return

	}
	//Plagin Name

	vrs := "v1"
	t.APIVersion = &vrs

	if *t.Module == "entrypoint" {
		errorAnswer(w, r, t, 401, "0000235", "Wrong module")
		return
	}

	if *t.Module == "heartbeat" {
		HeartbeatHandler(w, r, t)
		return
	}

	if r.Method == "POST" || r.Method == "DELETE" || r.Method == "PATCH" {
		t.Args = sf.ToMapStringAny(parseJSONArgs(r))
	}

	if r.Method == "GET" && r.URL.Query() != nil || r.Method == "TRACE" && r.URL.Query() != nil || r.Method == "HEAD" && r.URL.Query() != nil {
		paramMap := make(map[string]interface{}, 0)
		for k, v := range r.URL.Query() {
			if len(v) == 1 && len(v[0]) != 0 {
				paramMap[k] = v[0]
			}
		}
		anydt := sf.ToMapStringAny(paramMap)
		t.Args = anydt

	}

	//check for session
	if viper.GetBool("server.session") {
		t = checksession(t, r)

		if t.UID != nil && *t.Readonly == int32(1) {

			errorAnswer(w, r, t, 401, "0000235", "Read Only User")
			return

		}
	}

	//Load microservice
	if *t.Module == "info" {
		Info(w, r, t)
		return
	}
	connectgrpc(w, r, t)

}

func parseJSONArgs(r *http.Request) map[string]interface{} {
	var args map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		return nil
	}
	return args
}
