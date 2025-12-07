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
package registry

import (
	"errors"
	"sync"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

// --- Core cache structure ---
type ServiceInfo struct {
	Host       string
	Port       string
	LastUpdate time.Time
}

var (
	cache sync.Map // map[string]ServiceInfo
	ttl   = 60 * time.Second
)

// --- Public: get from cache or Master ---
func GetService(module string) (ServiceInfo, error) {
	if v, ok := cache.Load(module); ok {
		info := v.(ServiceInfo)
		if time.Since(info.LastUpdate) < ttl {
			return info, nil
		}
	}

	// Cache miss — try Master
	host := viper.GetString("microservices.masterservice.host")
	port := viper.GetString("microservices.masterservice.port")

	req := &pb.Request{
		Module: sf.StringPtr("masterservice"),
		IR: &pb.InternalRequest{
			Param:  sf.StringPtr("getmicroservicebypath"),
			Method: sf.StringPtr("GET"),
			Args:   map[string]*anypb.Any{},
		},
	}

	ans := sf.GRPCConnect(host, port, req)
	if ans["httpcode"] != nil {
		return ServiceInfo{}, errors.New("masterservice unavailable")
	}

	info := ServiceInfo{
		Host:       ans["host"].(string),
		Port:       ans["port"].(string),
		LastUpdate: time.Now(),
	}

	cache.Store(module, info)
	return info, nil
}

// --- Background refresh loop ---
func StartRefresher() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			RefreshCache()
		}
	}()
}

// --- Refresh all entries ---
func RefreshCache() {
	cache.Range(func(key, val any) bool {
		mod := key.(string)
		info := val.(ServiceInfo)

		if time.Since(info.LastUpdate) > ttl {
			// Lazy revalidation
			newInfo, err := GetService(mod)
			if err == nil {
				cache.Store(mod, newInfo)
			}
		}
		return true
	})
}

func StartSweeper() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			now := time.Now()
			cache.Range(func(key, value any) bool {
				info := value.(ServiceInfo)
				if now.Sub(info.LastUpdate) > ttl {
					cache.Delete(key)
				}
				return true
			})
		}
	}()
}
