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

func Error(code string) GufoError {
	if e, ok := Errors[code]; ok {
		return e
	}
	return GufoError{"99999", "Unknown Error", 500}
}

var Errors = map[string]GufoError{
	"00001": {"00001", "Unauthorized", 401},
	"00002": {"00002", "Invalid Session", 401},
	"00003": {"00003", "Bad Request", 400},
	"00004": {"00004", "Internal Error", 500},
}
