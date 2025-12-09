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

package registry

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

// ServiceInfo describes a resolved microservice endpoint.
type ServiceInfo struct {
	Host       string
	Port       string
	LastUpdate time.Time
}

var (
	cache sync.Map           // map[string]ServiceInfo
	ttl   = 60 * time.Second // default TTL for cached entries
)

// getRegistryMode returns normalized registry mode.
//
// Supported modes:
//   - "master"  – resolve via masterservice
//   - "static"  – resolve from local config/env
//
// If server.masterservice=true and no explicit mode set -> "master".
// Otherwise -> "static".
func getRegistryMode() string {
	mode := strings.ToLower(viper.GetString("server.registry_mode"))
	if mode != "" {
		return mode
	}

	if viper.GetBool("server.masterservice") {
		return "master"
	}
	return "static"
}

// getStaticServiceFromConfig resolves service from local config/env.
func getStaticServiceFromConfig(module string) (ServiceInfo, error) {
	keyPrefix := fmt.Sprintf("microservices.%s.", module)

	host := viper.GetString(keyPrefix + "host")
	port := viper.GetString(keyPrefix + "port")

	if host == "" || port == "" {
		return ServiceInfo{}, fmt.Errorf("static registry: microservice %q not found in config", module)
	}

	info := ServiceInfo{
		Host:       host,
		Port:       port,
		LastUpdate: time.Now(),
	}

	return info, nil
}

// getServiceFromMaster resolves service via masterservice microservice.
func getServiceFromMaster(module string) (ServiceInfo, error) {
	host := viper.GetString("microservices.masterservice.host")
	port := viper.GetString("microservices.masterservice.port")

	if host == "" || port == "" {
		return ServiceInfo{}, errors.New("masterservice host/port not configured")
	}

	req := &pb.Request{
		Module: sf.StringPtr("masterservice"),
		IR: &pb.InternalRequest{
			Param:  sf.StringPtr("getmicroservicebypath"),
			Method: sf.StringPtr("GET"),
			Args:   map[string]*anypb.Any{},
		},
	}

	ans := sf.GRPCConnect(host, port, req)
	if code, ok := ans["httpcode"]; ok && code != nil {
		return ServiceInfo{}, errors.New("masterservice unavailable")
	}

	info := ServiceInfo{
		Host:       ans["host"].(string),
		Port:       ans["port"].(string),
		LastUpdate: time.Now(),
	}

	return info, nil
}

// GetService resolves microservice endpoint either from cache,
// static config/env, or masterservice, depending on registry mode.
func GetService(module string) (ServiceInfo, error) {
	fmt.Fprintln(os.Stderr, ">>> REGISTRY GetService:", module)

	// 1️⃣ Cache
	if v, ok := cache.Load(module); ok {
		info := v.(ServiceInfo)
		if time.Since(info.LastUpdate) < ttl {
			return info, nil
		}
	}
	fmt.Fprintln(os.Stderr, ">>> REGISTRY cache miss:", module)

	// 2️⃣ STATIC REGISTRY MODE
	mode := strings.ToLower(viper.GetString("server.registry_mode"))
	fmt.Fprintln(os.Stderr, ">>> REGISTRY static lookup for:", module)

	if mode == "static" || viper.GetBool("server.masterservice") == false {

		key := strings.ReplaceAll(module, "-", "_")

		host := viper.GetString("microservices." + key + ".host")
		port := viper.GetString("microservices." + key + ".port")

		fmt.Fprintln(os.Stderr, ">>> REGISTRY static result host=", host, "port=", port)

		if host == "" || port == "" {
			return ServiceInfo{}, errors.New("static registry: host or port not set for " + key)
		}

		info := ServiceInfo{
			Host:       host,
			Port:       port,
			LastUpdate: time.Now(),
		}

		cache.Store(module, info)
		return info, nil
	}

	// 3️⃣ MASTER-SERVICE MODE
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
	fmt.Fprintln(os.Stderr, ">>> REGISTRY fallback to MASTER for:", module)

	ans := sf.GRPCConnect(host, port, req)
	if ans["httpcode"] != nil {
		fmt.Fprintln(os.Stderr, ">>> REGISTRY ERROR: cannot resolve", module)

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

// StartRefresher periodically revalidates cached entries.
// For "static" mode it is effectively a no-op, because config/env
// is treated as the source of truth.
func StartRefresher() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			RefreshCache()
		}
	}()
}

// RefreshCache updates expired entries depending on registry mode.
func RefreshCache() {
	mode := getRegistryMode()

	cache.Range(func(key, val any) bool {
		mod := key.(string)
		info := val.(ServiceInfo)

		if time.Since(info.LastUpdate) <= ttl {
			return true
		}

		var (
			newInfo ServiceInfo
			err     error
		)

		switch mode {
		case "static":
			// For static mode we simply reload from config/env.
			newInfo, err = getStaticServiceFromConfig(mod)
		case "master":
			newInfo, err = getServiceFromMaster(mod)
		default:
			err = fmt.Errorf("unknown registry mode: %s", mode)
		}

		if err == nil {
			cache.Store(mod, newInfo)
		}

		return true
	})
}

// StartSweeper removes expired entries from cache.
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
