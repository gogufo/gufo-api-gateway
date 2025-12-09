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
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
)

func JsonReq(url string, args map[string]interface{}, token string, tokentype string, tokenheader string, method string) ([]byte, error) {

	if method == "GET" {
		return JsonGet(url, args, token, tokentype, tokenheader)
	}

	json_data, err := json.Marshal(args)

	if err != nil {
		return nil, err
	}

	var jsonData = []byte(json_data)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if token != "" {
		header := token
		if tokentype != "" {
			header = tokentype + " " + token
		}
		req.Header.Add(tokenheader, header)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	byteresponse, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return byteresponse, nil
}
