// Copyright 2020-2024 Alexey Yanchenko <mail@yanchenko.me>
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

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

var methodHandlers = map[string]func(http.ResponseWriter, *http.Request, *pb.Request, int){
	"OPTIONS": ProcessOPTIONS,
	"GET":     ProcessREQ,
	"HEAD":    ProcessREQ,
	"TRACE":   ProcessREQ,
	"POST":    ProcessREQ,
	"PATCH":   ProcessREQ,
	"DELETE":  ProcessREQ,
	"PUT":     ProcessPUT,
}

func API(w http.ResponseWriter, r *http.Request, version int) {
	t := RequestInit(r)

	if h, ok := methodHandlers[r.Method]; ok {
		h(w, r, t, version)
	} else {
		ProcessOPTIONS(w, r, t, version)
	}

}
