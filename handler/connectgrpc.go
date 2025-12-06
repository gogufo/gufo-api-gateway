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
	"github.com/gogufo/gufo-api-gateway/registry"
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

	// ------------------------------------------------------------
	// CLUSTER MODE: resolve via MasterService
	// ------------------------------------------------------------
	if *t.Module != "masterservice" && msmethod {

		host = viper.GetString("microservices.masterservice.host")
		port = viper.GetString("microservices.masterservice.port")

		// Backup original InternalRequest
		origIR := t.IR
		defer func() { t.IR = origIR }()

		mst := &pb.InternalRequest{}
		param := "getmicroservicebypath"
		gt := "GET"
		mst.Param = &param
		mst.Method = &gt
		t.IR = mst

		// ---- TIMEOUT WRAPPER WITHOUT NEW API ----
		ansChan := make(chan map[string]interface{}, 1)
		go func() {
			ansChan <- sf.GRPCConnect(host, port, t)
		}()

		var ans map[string]interface{}
		select {
		case ans = <-ansChan:
		case <-time.After(3 * time.Second):
			ans = map[string]interface{}{"httpcode": 504}
		}

		if ans == nil || ans["httpcode"] != nil {

			// MasterService unavailable → local registry fallback
			cached, err := registry.GetService(*t.Module)
			if err == nil {
				sf.SetLog(fmt.Sprintf("⚠️ Using cached route for %s (%s:%s)", *t.Module, cached.Host, cached.Port))
				return cached.Host, cached.Port, ""
			}

			msg := fmt.Sprintf("MasterService unavailable and no cached entry for %s", *t.Module)
			if viper.GetBool("server.sentry") {
				sentry.CaptureMessage(msg)
			} else {
				sf.SetErrorLog(msg)
			}
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

		return host, port, plygintype
	}

	// ------------------------------------------------------------
	// STANDALONE MODE: resolve from config
	// ------------------------------------------------------------
	if !viper.IsSet(pluginname) {
		msg := fmt.Sprintf("No Module %s", *t.Module)
		if viper.GetBool("server.sentry") {
			sentry.CaptureMessage(msg)
		} else {
			sf.SetErrorLog(msg)
		}
		return "", "", ""
	}

	hostpath := fmt.Sprintf("%s.host", pluginname)
	portpath := fmt.Sprintf("%s.port", pluginname)

	host = viper.GetString(hostpath)
	port = viper.GetString(portpath)

	plygintype = viper.GetString(fmt.Sprintf("%s.type", pluginname))

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

	// ------------------------------------------------------------
	// Resolve service from registry, fallback to GetHostAndPort
	// ------------------------------------------------------------
	info, err := registry.GetService(*t.Module)
	if err != nil {

		host, port, _ := GetHostAndPort(t)
		if host == "" || port == "" {
			errorAnswer(w, r, t, 500, "0000501", "Cannot resolve service: registry and masterservice unavailable")
			return
		}

		// Create local info object WITHOUT touching registry internals
		info = registry.ServiceInfo{
			Host: host,
			Port: port,
		}
	}

	// ------------------------------------------------------------
	// 2️⃣ Streaming uploads (PUT)
	// ------------------------------------------------------------
	if r.Method == http.MethodPut {
		ans := sf.GRPCStreamPut(info.Host, info.Port, r, t)
		moduleAnswerv3(w, r, ans, t)
		return
	}

	// ------------------------------------------------------------
	// 3️⃣ Standard transport call
	// ------------------------------------------------------------
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
