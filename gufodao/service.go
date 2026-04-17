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
	"net/http"
	"strings"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

/*
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
*/
// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func ReadUserIP(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

// StringPtr returns a pointer to a string value.
func StringPtr(v string) *string {
	return &v
}

// Int32Ptr returns a pointer to an int32 value.
func Int32Ptr(v int32) *int32 {
	return &v
}

func ProtoMethodToString(m pb.Method) string {
	switch m {
	case pb.Method_METHOD_GET:
		return http.MethodGet
	case pb.Method_METHOD_POST:
		return http.MethodPost
	case pb.Method_METHOD_PUT:
		return http.MethodPut
	case pb.Method_METHOD_PATCH:
		return http.MethodPatch
	case pb.Method_METHOD_DELETE:
		return http.MethodDelete
	default:
		return ""
	}
}

// Convert HTTP method string to proto enum
func HttpMethodToProto(m string) pb.Method {
	switch m {
	case http.MethodGet:
		return pb.Method_METHOD_GET
	case http.MethodPost:
		return pb.Method_METHOD_POST
	case http.MethodPut:
		return pb.Method_METHOD_PUT
	case http.MethodPatch:
		return pb.Method_METHOD_PATCH
	case http.MethodDelete:
		return pb.Method_METHOD_DELETE
	default:
		return pb.Method_METHOD_UNSPECIFIED
	}
}
