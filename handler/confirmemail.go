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

type ConfEmailLink struct {
	Email string
	Token string
	Lang  string
}

/*
func Confirmemail(w http.ResponseWriter, r *http.Request) {

	t := RequestInit(r)
	//Log Request
	//1. Collect need data
	var userip = sf.ReadUserIP(r)
	sf.SetLog(userip + " /confirmemail " + r.Method)

	ans := make(map[string]interface{})
	p := bluemonday.UGCPolicy()

	//1. We get request with email and hash and check if it need data exist
	if r.URL.Query()["token"][0] == "" || r.URL.Query()["email"][0] == "" {
		errorAnswer(w, r, t, 400, "000001", "Missing Token or Email")
		return
	}

	//2. Clen it from any tags
	var data ConfEmailLink
	data.Email = p.Sanitize(r.URL.Query()["email"][0])
	data.Token = p.Sanitize(r.URL.Query()["token"][0])
	if r.URL.Query()["lang"][0] == "" {
		data.Lang = "en"
	} else {
		data.Lang = p.Sanitize(r.URL.Query()["lang"][0])
	}

	//4. Check is hash live
	var userHash sf.TimeHash

	//Check DB and table config
	db, err := sf.ConnectDBv2()
	if err != nil {

		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("confirmemail.go: " + err.Error())
		}
		//return "error with db"
		errorAnswer(w, r, t, 400, "000001", "DB Connection Error")
	}

	//4.1. Check if hash is exist in db users
	rows := db.Conn.Where(`hash = ? and mail = ?`, data.Token, data.Email).First(&userHash)
	if rows.RowsAffected == 0 {
		// return error. Hash is not exist in db
		errorAnswer(w, r, t, 400, "000008", "Hash is not exist in db")
		return
	}

	curtime := int(time.Now().Unix())
	if userHash.Livetime < curtime {
		errorAnswer(w, r, t, 400, "000009", "Hash is overtime")
		return
	}

	//5. Update users table
	db.Conn.Table("users").Where("mail = ?", data.Email).Updates(map[string]interface{}{"completed": true, "mailconfirmed": curtime})

	//6. Delete hash
	db.Conn.Delete(sf.TimeHash{}, "hash = ?", data.Token)

	ans["response"] = "100002" // email confirmed
	moduleAnswerv3(w, r, ans, t)
	return

}
*/
