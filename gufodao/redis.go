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

import (
	"crypto/tls"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var CachePool *redis.Pool

func InitCache() {
	host := ConfigString("redis.host")
	password := ConfigString("redis.password")
	useTLS := ConfigBool("redis.tls")
	maxIdle := viper.GetInt("redis.max_idle")
	if maxIdle == 0 {
		maxIdle = 5
	}

	maxActive := viper.GetInt("redis.max_active")
	if maxActive == 0 {
		maxActive = 20
	}

	idleTimeout := viper.GetInt("redis.idle_timeout")
	if idleTimeout == 0 {
		idleTimeout = 240
	}

	var tlsConfig *tls.Config
	if useTLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	CachePool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Wait:        true,

		Dial: func() (redis.Conn, error) {
			options := []redis.DialOption{
				redis.DialConnectTimeout(3 * time.Second),
				redis.DialReadTimeout(3 * time.Second),
				redis.DialWriteTimeout(3 * time.Second),
			}

			if password != "" {
				options = append(options, redis.DialPassword(password))
			}

			if tlsConfig != nil {
				options = append(options, redis.DialTLSConfig(tlsConfig))
			}

			conn, err := redis.Dial("tcp", host, options...)
			if err != nil {
				SetErrorLog("redis: dial failed: " + err.Error())
				return nil, err
			}
			return conn, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 30*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				SetErrorLog("redis: ping failed: " + err.Error())
			}
			return err
		},
	}

	// initial test
	conn := CachePool.Get()
	_, err := conn.Do("PING")
	conn.Close()

	if err != nil {
		SetErrorLog("redis: initial ping failed: " + err.Error())
	} else {
		SetLog("redis: connected")
	}
}
