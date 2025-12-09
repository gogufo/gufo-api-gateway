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

package handler

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/spf13/viper"
)

func moduleAnswerv3(w http.ResponseWriter, r *http.Request, s map[string]interface{}, t *pb.Request) {

	// 🧩 Make a copy of the incoming map to avoid side effects
	out := make(map[string]interface{})
	for k, v := range s {
		out[k] = v
	}

	// --- File response handler ---
	if out["file"] != nil {
		filename := out["file"].(string)
		base64type := false
		if out["isbase64"] != nil {
			base64type = out["isbase64"].(bool)
		}
		fileAnswer(w, r, filename, out["filetype"].(string), out["filename"].(string), base64type)
		return
	}

	// --- Standard JSON response ---
	var resp sf.Response
	httpsstatus := 200

	// Extract HTTP code from map (if provided)
	if out["httpcode"] != nil {
		switch reflect.TypeOf(out["httpcode"]).String() {
		case "string":
			pre := out["httpcode"].(string)
			httpsstatus, _ = strconv.Atoi(pre)
		case "int":
			httpsstatus = out["httpcode"].(int)
		case "float64":
			httpsstatus = int(out["httpcode"].(float64))
		}
		delete(out, "httpcode")
	}

	// Allow microservice to override Content-Type
	if ct, ok := out["Content-Type"]; ok {
		if cts, ok2 := ct.(string); ok2 {
			w.Header().Set("Content-Type", cts)
			delete(out, "Content-Type")
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	// Propagate X-Request-ID if present
	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		w.Header().Set("X-Request-ID", rid)
	}

	// Language field
	resp.Language = "eng"
	if out["lang"] != nil {
		resp.Language = out["lang"].(string)
	}

	// Timestamp
	resp.TimeStamp = int(time.Now().Unix())
	resp.Data = out

	// Attach session info if present
	if t.UID != nil {
		session := make(map[string]interface{})
		session["uid"] = t.UID
		session["isAdmin"] = t.IsAdmin
		session["Sesionexp"] = t.SessionEnd
		session["completed"] = t.Completed
		session["readonly"] = t.Readonly
		resp.Session = session
	}

	// Marshal response to JSON
	answer, err := json.Marshal(resp)
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("api.go: " + err.Error())
		}
		return
	}

	// Apply default headers
	for i := 0; i < len(HeaderKeys); i++ {
		w.Header().Set(HeaderKeys[i], HeaderValues[i])
	}

	// Final response
	w.WriteHeader(httpsstatus)
	w.Write([]byte(answer))
}
