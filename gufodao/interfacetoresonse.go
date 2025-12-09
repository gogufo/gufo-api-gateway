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
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Interfacetoresponse(request *pb.Request, answer map[string]interface{}) (response *pb.Response) {

	decanswer := ToMapStringAny(answer)
	response = &pb.Response{
		Data:        decanswer,
		RequestBack: request,
	}

	return response

}

func ErrorReturn(t *pb.Request, httpcode int, code string, message string) (response *pb.Response) {

	ans := make(map[string]interface{})

	ans["httpcode"] = httpcode
	ans["code"] = code
	ans["message"] = message
	response = Interfacetoresponse(t, ans)

	return response
}
