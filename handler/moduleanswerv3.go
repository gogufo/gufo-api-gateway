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
	"encoding/json"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufodao"

	"github.com/spf13/viper"
)

func moduleAnswerv3(w http.ResponseWriter, r *http.Request, s map[string]interface{}, t *sf.Request) {

	if s["file"] != nil {
		var filename = s["file"].(string)
		fileAnswer(w, r, filename, s["filetype"].(string), s["filename"].(string))
	} else {
		var resp sf.Response

		httpsstatus := 200

		if s["httpcode"] != nil {
			httpsstatus = s["httpcode"].(int)
		}

		resp.Language = "eng"

		if s["lang"] != nil {
			resp.Language = s["lang"].(string)
		}

		resp.TimeStamp = int(time.Now().Unix())
		resp.Data = s

		if t.UID != "" {
			//write session data in answer
			session := make(map[string]interface{})
			session["uid"] = t.UID
			session["isAdmin"] = t.IsAdmin
			session["Sesionexp"] = t.SessionEnd
			session["completed"] = t.Completed
			session["readonly"] = t.Readonly
			resp.Session = session
		}

		answer, err := json.Marshal(resp)
		if err != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog("api.go: " + err.Error())
			}
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Server", "Gufo")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpsstatus)
		w.Write([]byte(answer))
	}

}
