// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
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
