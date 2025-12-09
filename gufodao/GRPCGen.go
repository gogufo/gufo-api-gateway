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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	viper "github.com/spf13/viper"
)

func GRPCGen(misroservice string, param string, paramid string, args map[string]interface{}, token string, method string, sign string) map[string]interface{} {

	ans := make(map[string]interface{})

	erphost := viper.GetString("server.internal_host")
	erpport := viper.GetString("server.port")
	tsp := "http://"
	isssl := viper.GetBool("server.internal_ssl")
	if isssl {
		tsp = "https://"
	}

	header := "Bearer " + token
	URL := fmt.Sprintf("%s%s:%s/api/v3/%s/%s", tsp, erphost, erpport, misroservice, param)
	if paramid != "" {
		URL = fmt.Sprintf("%s/%s", URL, paramid)
	}

	if len(args) != 0 {

		var b []string
		for key, value := range args {
			str := fmt.Sprintf("%s=%s", key, value)
			b = append(b, str)
		}
		URLValues := strings.Join(b, "&")
		URL = fmt.Sprintf("%s?%s", URL, URLValues)

	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		ans["error"] = err.Error()
		ans["httpcode"] = 400
		//	return ErrorReturn(t, 400, "000005", err.Error())

	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {header},
		"X-Sign":        {sign},
	}

	res, err := client.Do(req)
	if err != nil {
		ans["error"] = err.Error()
		ans["httpcode"] = 400
		//return ErrorReturn(t, 400, "000005", err.Error())
	}

	var cResp Response

	if err = json.NewDecoder(res.Body).Decode(&cResp); err != nil {
		//	return ErrorReturn(t, 400, "000005", err.Error())
		ans["error"] = err.Error()
		ans["httpcode"] = 400
	}

	ans["answer"] = cResp
	ans["httpcode"] = res.StatusCode

	return ans
}
