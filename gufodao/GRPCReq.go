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
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	viper "github.com/spf13/viper"
)

// HttpRetry is number of HTTP retries for internal requests
const HttpRetry = 3

// HttpTimeout is timeout per request
const HttpTimeout = 7 * time.Second

// Production-ready internal API Gateway request
func GRPCReq(
	microservice string,
	param string,
	paramID string,
	args map[string]interface{},
	token string,
	method string,
	sign string,
) map[string]interface{} {

	ans := make(map[string]interface{})

	// ---------------------------
	// Build URL
	// ---------------------------
	host := viper.GetString("server.internal_host")
	port := viper.GetString("server.port")

	if host == "" || port == "" {
		ans["error"] = "internal_host or port is empty"
		ans["httpcode"] = 500
		return ans
	}

	proto := "http://"
	if viper.GetBool("server.internal_ssl") {
		proto = "https://"
	}

	url := fmt.Sprintf("%s%s:%s/api/v3/%s/%s", proto, host, port, microservice, param)
	if paramID != "" {
		url += "/" + paramID
	}

	// ---------------------------
	// Prepare JSON body
	// ---------------------------
	var body []byte
	if args != nil {
		var err error
		body, err = json.Marshal(args)
		if err != nil {
			ans["error"] = "json marshal error: " + err.Error()
			ans["httpcode"] = 400
			return ans
		}
	} else {
		body = []byte("{}")
	}

	// ---------------------------
	// Prepare HTTP client
	// ---------------------------
	transport := &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
		TLSHandshakeTimeout: 5 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	// secure TLS
	if viper.GetBool("server.internal_ssl") {
		transport.TLSClientConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false, // production: must be false
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   HttpTimeout,
	}

	// ---------------------------
	// RETRY logic
	// ---------------------------
	var lastErr error
	for attempt := 1; attempt <= HttpRetry; attempt++ {

		// context per request
		ctx, cancel := context.WithTimeout(context.Background(), HttpTimeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
		if err != nil {
			lastErr = err
			continue
		}

		// headers
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Gufo-Microservice/2.0")
		req.Header.Set("X-Sign", sign)

		// execute
		res, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt) * 150 * time.Millisecond)
			continue
		}

		// ensure body close
		defer res.Body.Close()

		// parse response
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			lastErr = err
			continue
		}

		var cResp Response
		if err := json.Unmarshal(bodyBytes, &cResp); err != nil {
			lastErr = err
			continue
		}

		// success
		ans["answer"] = cResp
		ans["httpcode"] = res.StatusCode
		return ans
	}

	// if reached — all attempts failed
	ans["error"] = fmt.Sprintf("request failed after %d attempts: %v", HttpRetry, lastErr)
	ans["httpcode"] = 500
	return ans
}
