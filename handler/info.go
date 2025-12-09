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
	"fmt"
	"net/http"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	v "github.com/gogufo/gufo-api-gateway/version"

	"github.com/spf13/viper"
)

func Info(w http.ResponseWriter, r *http.Request, t *pb.Request) {
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /info " + r.Method)

	msg := fmt.Sprintf("%s  (%s, %s)", v.VERSION, v.GitCommit, v.BuildDate)

	ans := make(map[string]interface{})
	ans["version"] = msg
	ans["registration"] = viper.GetBool("settings.registration")

	moduleAnswerv3(w, r, ans, t)

}
