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
// Copyright 2025 Alexey Yanchenko <mail@yanchenko.me>
// SPDX-License-Identifier: Apache-2.0

package gufodao

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger      *zap.Logger
	errorLogger *zap.Logger
	once        sync.Once
)

// InitLogger initializes the JSON logger with daily rotation and size limit.
func InitLogger() {
	once.Do(func() {
		logDir := GetLogDir()
		if logDir == "" {
			logDir = "/var/gufo/log/"
		}

		// Create directories if needed
		_ = EnsureDir(logDir)

		// Daily log file name (e.g. gufo-2025-11-01.log)
		date := time.Now().Format("2006-01-02")
		mainLogPath := filepath.Join(logDir, "gufo-"+date+".log")
		errorLogPath := filepath.Join(logDir, "error-"+date+".log")

		mainWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   mainLogPath,
			MaxSize:    10, // MB
			MaxBackups: 7,
			MaxAge:     30, // days
			Compress:   true,
		})

		errorWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   errorLogPath,
			MaxSize:    10,
			MaxBackups: 7,
			MaxAge:     30,
			Compress:   true,
		})

		encoderCfg := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			mainWriter,
			zapcore.InfoLevel,
		)

		errCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			errorWriter,
			zapcore.ErrorLevel,
		)

		logger = zap.New(core, zap.AddCaller())
		errorLogger = zap.New(errCore, zap.AddCaller())
	})
}

// EnsureDir makes sure log directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// SetLog writes an info message to the main log in JSON format
func SetLog(msg string) {
	if logger == nil {
		InitLogger()
	}
	logger.Info(msg)
}

// SetErrorLog writes an error message to the error log in JSON format
func SetErrorLog(msg string) {
	if errorLogger == nil {
		InitLogger()
	}
	errorLogger.Error(msg)
}

// FlushLog ensures buffered data is written to disk (call before exit)
func FlushLog() {
	if logger != nil {
		_ = logger.Sync()
	}
	if errorLogger != nil {
		_ = errorLogger.Sync()
	}
}

func GetLogDir() string {
	n := ConfigString("server.logdir")
	var logdir string = n
	return logdir
}
