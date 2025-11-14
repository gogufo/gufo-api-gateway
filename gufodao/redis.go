// Copyright 2019 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Neptune library.
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
