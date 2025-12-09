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
//
// This file content Handler for API
// Each API function is independend plugin
// and API get reguest in connect with plugin
// Get response from plugin and answer to client
// All data is in JSON format

package handler

import (
	"bytes"
	"net/http"
	"time"
)

func fileAnswer(w http.ResponseWriter, r *http.Request, filepath string, filetype string, filename string, base64type bool) {
	for i := 0; i < len(HeaderKeys); i++ {
		if HeaderKeys[i] != "Content-Type" {
			w.Header().Set(HeaderKeys[i], HeaderValues[i])
		}
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", filetype)
	w.WriteHeader(200)
	if base64type {
		http.ServeContent(w, r, filename, time.Time{}, bytes.NewReader([]byte(filepath))) //ServerContent for base64 files
	} else {
		http.ServeFile(w, r, filepath) //ServerFile for download files
	}

}
