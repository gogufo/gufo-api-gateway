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

package handler

import (
	"net/http"
	"strings"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func RequestInit(r *http.Request) *pb.Request {
	t := &pb.Request{}
	p := bluemonday.UGCPolicy()

	path := r.URL.Path
	patharray := strings.Split(path, "/")
	pathlenth := len(patharray)

	module := p.Sanitize(patharray[3])
	t.Module = &module
	t.Path = &path
	t.Method = &r.Method

	sgn := viper.GetString("server.sign")
	curip := sf.ReadUserIP(r)
	usagent := r.UserAgent()

	t.Sign = &sgn

	t.IP = &curip

	t.UserAgent = &usagent

	//Function in Plugin
	if pathlenth >= 5 {
		ptr := p.Sanitize(patharray[4])
		t.Param = &ptr
	}

	//ID for function in plugin
	if pathlenth >= 6 {
		ptrs := p.Sanitize(patharray[5])
		t.ParamID = &ptrs

	}

	if pathlenth >= 7 {
		ptrs := p.Sanitize(patharray[6])
		t.ParamIDD = &ptrs

	}

	return t
}
