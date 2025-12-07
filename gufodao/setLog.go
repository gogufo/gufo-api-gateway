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
package gufodao

import (
	"os"
	"path/filepath"
	"strings"
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

// InitLogger initializes the JSON logger.
// If GUFO_LOG_TO_FILE=false -> logs go to stdout/stderr (Kubernetes-friendly).
// If GUFO_LOG_TO_FILE=true  -> logs are written to rotating files.
func InitLogger() {
	once.Do(func() {
		logToFile := strings.ToLower(os.Getenv("GUFO_LOG_TO_FILE")) == "true"

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

		var mainWriter zapcore.WriteSyncer
		var errorWriter zapcore.WriteSyncer

		if logToFile {
			logDir := GetLogDir()
			if logDir == "" {
				logDir = "/var/gufo/log/"
			}

			_ = EnsureDir(logDir)

			date := time.Now().Format("2006-01-02")
			mainLogPath := filepath.Join(logDir, "gufo-"+date+".log")
			errorLogPath := filepath.Join(logDir, "error-"+date+".log")

			mainWriter = zapcore.AddSync(&lumberjack.Logger{
				Filename:   mainLogPath,
				MaxSize:    10, // MB
				MaxBackups: 7,
				MaxAge:     30, // days
				Compress:   true,
			})

			errorWriter = zapcore.AddSync(&lumberjack.Logger{
				Filename:   errorLogPath,
				MaxSize:    10,
				MaxBackups: 7,
				MaxAge:     30,
				Compress:   true,
			})
		} else {
			// ✅ Kubernetes / container-friendly logging
			mainWriter = zapcore.AddSync(os.Stdout)
			errorWriter = zapcore.AddSync(os.Stderr)
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

		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		errorLogger = zap.New(errCore, zap.AddCaller(), zap.AddCallerSkip(1))
	})
}

// EnsureDir makes sure log directory exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// SetLog writes an info message
func SetLog(msg string) {
	if logger == nil {
		InitLogger()
	}
	logger.Info(msg)
}

// SetErrorLog writes an error message
func SetErrorLog(msg string) {
	if errorLogger == nil {
		InitLogger()
	}
	errorLogger.Error(msg)
}

// FlushLog ensures buffered data is written (noop for stdout)
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
