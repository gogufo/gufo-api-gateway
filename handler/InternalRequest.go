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
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

/*
JSON /api/v3/auth/signin POST username and password
GRPC /session/savesession POST
*/

func InternalRequest(t *pb.Request) (response *pb.Response) {
	//Get destination way
	if t.Module != nil && *t.Module == "heartbeat" {
		// тут payload можно собрать из t.Args, если нужно
		ans, err := heartbeatCore(t, nil)
		if err != nil {
			return sf.ErrorReturn(t, 500, "0000501", "MasterService heartbeat error")
		}
		return sf.Interfacetoresponse(t, ans)
	}

	host, port, _ := GetHostAndPort(t)

	ans := sf.GRPCConnect(host, port, t)

	response = sf.Interfacetoresponse(t, ans)

	return response
}
