// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package handler

import (
	"net/http"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func WrongRequest(w http.ResponseWriter, r *http.Request) {
	//Log Request
	//1. Collect need data

	//check for session
	ans := make(map[string]interface{})
	ans["code"] = "000056"
	ans["mesage"] = "Wrong Request"
	ans["httpcode"] = 404

	t := &pb.Request{}

	moduleAnswerv3(w, r, ans, t)

}
