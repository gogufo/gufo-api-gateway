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
// This file contains the main response handler for Gufo API Gateway.
// Each API module acts independently and returns JSON-formatted data.

// Copyright 2019-2025 Alexey Yanchenko <mail@yanchenko.me>

package handler

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

func moduleAnswerv3(w http.ResponseWriter, r *http.Request, s map[string]interface{}, t *pb.Request) {

	// --- Copy input to avoid mutation ---
	out := make(map[string]interface{}, len(s))
	for k, v := range s {
		out[k] = v
	}

	// ------------------------------------------------------------
	// FILE RESPONSE (bypass JSON)
	// ------------------------------------------------------------
	if file, ok := out["file"]; ok {
		filename := file.(string)

		base64type := false
		if v, ok := out["isbase64"]; ok {
			base64type = v.(bool)
		}

		fileAnswer(
			w,
			r,
			filename,
			out["filetype"].(string),
			out["filename"].(string),
			base64type,
		)
		return
	}

	// ------------------------------------------------------------
	// HTTP STATUS EXTRACTION (legacy support)
	// ------------------------------------------------------------
	httpsstatus := 200

	if v, ok := out["httpcode"]; ok {
		switch reflect.TypeOf(v).String() {
		case "string":
			httpsstatus, _ = strconv.Atoi(v.(string))
		case "int":
			httpsstatus = v.(int)
		case "float64":
			httpsstatus = int(v.(float64))
		}
		delete(out, "httpcode")
	}

	// ------------------------------------------------------------
	// HEADERS
	// ------------------------------------------------------------
	if ct, ok := out["Content-Type"]; ok {
		if cts, ok2 := ct.(string); ok2 {
			w.Header().Set("Content-Type", cts)
		}
		delete(out, "Content-Type")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		w.Header().Set("X-Request-ID", rid)
	}

	// ------------------------------------------------------------
	// RESPONSE STRUCT (legacy JSON envelope)
	// ------------------------------------------------------------
	var resp sf.Response
	resp.Data = out

	// ------------------------------------------------------------
	// SESSION (temporary backward compatibility)
	// ------------------------------------------------------------
	if t.Auth != nil && t.Auth.Uid != "" {

		session := map[string]interface{}{
			"uid":       t.Auth.Uid,
			"isAdmin":   t.Auth.IsAdmin,
			"Sesionexp": t.Auth.SessionEnd,
			"readonly":  t.Auth.Readonly,
		}

		// optional meta fields
		if t.Auth.Meta != nil {
			if v, ok := t.Auth.Meta["completed"]; ok {
				session["completed"] = v
			}
		}

		resp.Session = session
	}

	// ------------------------------------------------------------
	// JSON SERIALIZATION
	// ------------------------------------------------------------
	answer, err := json.Marshal(resp)
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("moduleAnswerv3: " + err.Error())
		}
		return
	}

	// ------------------------------------------------------------
	// DEFAULT HEADERS
	// ------------------------------------------------------------
	for i := 0; i < len(HeaderKeys); i++ {
		w.Header().Set(HeaderKeys[i], HeaderValues[i])
	}

	// ------------------------------------------------------------
	// FINAL RESPONSE
	// ------------------------------------------------------------
	w.WriteHeader(httpsstatus)
	w.Write(answer)
}
