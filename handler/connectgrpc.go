// Copyright 2024 Alexey Yanchenko <mail@yanchenko.me>
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

package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/gogufo/gufo-api-gateway/transport"
	"github.com/spf13/viper"
	pbv "gopkg.in/cheggaaa/pb.v1"
)

const chunkSize = 64 * 1024

type uploader struct {
	ctx         context.Context
	wg          sync.WaitGroup
	requests    chan string // each request is a filepath on client accessible to client
	pool        *pbv.Pool
	DoneRequest chan string
	FailRequest chan string
}

func GetHostAndPort(t *pb.Request) (host string, port string, plygintype string) {

	pluginname := fmt.Sprintf("microservices.%s", *t.Module)

	msmethod := viper.GetBool("server.masterservice")

	if *t.Module != "masterservice" && msmethod {

		//Check masterservice for host and port
		host = viper.GetString("microservices.masterservice.host")
		port = viper.GetString("microservices.masterservice.port")

		//Modify data for request masterservice
		mst := &pb.InternalRequest{}
		param := "getmicroservicebypath"
		gt := "GET"
		mst.Param = &param
		mst.Method = &gt

		t.IR = mst

		ans := sf.GRPCConnect(host, port, t)
		if ans["httpcode"] != nil {

			if viper.GetBool("server.sentry") {
				sentry.CaptureMessage(fmt.Sprintf("%v", ans["message"]))
			} else {
				sf.SetErrorLog(fmt.Sprintf("%v", ans["message"]))
			}
			//	httpcode := 0

			//	httpcode, _ = strconv.Atoi(fmt.Sprintf("%v", ans["httpcode"]))

			//	errorAnswer(w, r, t, httpcode, fmt.Sprintf("%v", ans["code"]), fmt.Sprintf("%v", ans["message"]))
			return "", "", ""
		}

		host = fmt.Sprintf("%v", ans["host"])
		port = fmt.Sprintf("%v", ans["port"])
		if ans["isinternal"] != nil {
			isint, _ := strconv.ParseBool(fmt.Sprintf("%v", ans["isinternal"]))
			if isint {
				plygintype = "internal"
			}
		}

		//Put previoud data back
		//	*st.Request.Param = curparam
		//	*st.Request.Method = curmethod

	} else {

		if !viper.IsSet(pluginname) {
			msg := fmt.Sprintf("No Module %s", *t.Module)
			if viper.GetBool("server.sentry") {
				sentry.CaptureMessage(msg)
			} else {
				sf.SetErrorLog(msg)
			}
			//	errorAnswer(w, r, t, 401, "0000235", msg)
			return "", "", ""
		}

		hostpath := fmt.Sprintf("%s.host", pluginname)
		portpath := fmt.Sprintf("%s.port", pluginname)
		host = viper.GetString(hostpath)
		port = viper.GetString(portpath)
	}

	if !msmethod {
		plygintype = fmt.Sprintf("%s.type", pluginname)
	}

	return host, port, plygintype
}

func connectgrpc(w http.ResponseWriter, r *http.Request, t *pb.Request) {
	// 1️⃣ Internal signature check (optional)
	if r.Header.Get("X-Sign") != "" {
		sign := r.Header.Get("X-Sign")
		expected := viper.GetString("server.sign")
		if sign != expected {
			errorAnswer(w, r, t, 401, "0000234", "Invalid internal signature")
			return
		}
	}

	// 2️⃣ Handle streaming uploads (PUT)
	if r.Method == http.MethodPut {
		host := viper.GetString(fmt.Sprintf("microservices.%s.host", *t.Module))
		port := viper.GetString(fmt.Sprintf("microservices.%s.port", *t.Module))

		ans := sf.GRPCStreamPut(host, port, r, t)
		moduleAnswerv3(w, r, ans, t)
		return
	}

	// 3️⃣ Standard call via transport abstraction
	tr := transport.Get()
	ctx := r.Context()

	resp, err := tr.Call(ctx, *t.Module, r.Method, t)
	if err != nil {
		errorAnswer(w, r, t, 500, "0000500", err.Error())
		return
	}

	moduleAnswerv3(w, r, sf.ToMapStringInterface(resp.Data), t)
}

func (d *uploader) Stop() {
	close(d.requests)
	d.wg.Wait()
	d.pool.RefreshRate = 500 * time.Millisecond
	d.pool.Stop()
}
