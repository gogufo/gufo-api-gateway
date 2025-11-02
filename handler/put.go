// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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

	pb "github.com/gogufo/gufo-api-gateway/proto/go"

	"github.com/spf13/viper"
)

func ProcessPUT(w http.ResponseWriter, r *http.Request, t *pb.Request, version int) {
	path := r.URL.Path
	patharray := strings.Split(path, "/")
	if len(patharray) < 3 || *t.Module == "entrypoint" {
		errorAnswer(w, r, t, 401, "0000235", "Wrong Path Length")
		return
	}

	// 🔐 Check session
	if viper.GetBool("server.session") {
		t = checksession(t, r)
		if t.UID != nil && *t.Readonly == int32(1) {
			errorAnswer(w, r, t, 401, "0000235", "Read Only User")
			return
		}
	}

	vrs := "v3"
	if version == 2 {
		vrs = "v2"
	}
	t.APIVersion = &vrs

	param := "stream"
	t.IR = &pb.InternalRequest{Param: &param}

	connectgrpc(w, r, t)
}
