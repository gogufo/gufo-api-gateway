// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
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

package gufodao

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	viper "github.com/spf13/viper"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	configDir  = "config"
	configName = "settings"
)

// defaultConfigExample is a minimal safe configuration
// automatically written when no config file is found.
var defaultConfigExample = []byte(`# Auto-generated default config (safe defaults)
[server]
port = "8090"
grpc_port = "4890"
debug = false
ip = "0.0.0.0"
sentry = false
session = true
masterservice = true
sysdir = "/var/gufo/"
tempdir = "/var/gufo/templates/"
filedir = "/var/gufo/files/"
plugindir = "/var/gufo/lib/"
logdir = "/var/gufo/log/"

[security]
# Sensitive values should be provided via ENV (see .env.example):
# GUFO_SIGN, GUFO_JWT_SECRET
sign_env = "GUFO_SIGN"
jwt_secret_env = "GUFO_JWT_SECRET"

[microservices.masterservice]
host = "masterservice"
port = "5300"
type = "server"
entrypointversion = "1.0.0"
cron = false

[database]
type = "mysql"
host = "db"
port = "3306"
dbname = "gufo"
user = "root"
password_env = "GUFO_DB_PASS"
charset = "utf8mb4"
protocol = "tcp"
sslmode = "disable"

[redis]
host = "redis://redis"

[sentry]
enabled = false
dsn_env = "GUFO_SENTRY_DSN"
trace = 1.0
flush = 2
tracing = true
debug = false
`)

// EnsureConfigExists verifies if config/settings.toml exists.
// If it does not, the function creates it from an embedded safe default.
// No interactive prompts or secrets are written.
func EnsureConfigExists() {
	path := filepath.Join(configDir, configName+".toml")

	info, statErr := os.Stat(path)
	if statErr == nil && !info.IsDir() {
		return // file already exists
	}

	if statErr != nil && !os.IsNotExist(statErr) {
		SetErrorLog("config: cannot stat config file: " + statErr.Error())
		return
	}

	if mkErr := os.MkdirAll(configDir, 0o755); mkErr != nil {
		SetErrorLog("config: cannot create config dir: " + mkErr.Error())
		return
	}

	if writeErr := os.WriteFile(path, defaultConfigExample, 0o644); writeErr != nil {
		SetErrorLog("config: cannot write default config: " + writeErr.Error())
		return
	}

	SetLog("config: created default config at " + path)
}

// InitConfig loads configuration with layered fallback.
// It loads .env (optional), reads environment overrides,
// parses TOML file if present, applies defaults, and validates critical keys.
func InitConfig() error {
	// 1) Load .env for local/non-Docker usage
	_ = godotenv.Load()

	// 2) Configure Viper for ENV overrides
	viper.SetEnvPrefix("GUFO")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 3) Set config search paths
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".") // fallback when run from project root

	// 4) Safe defaults to allow startup without config file
	viper.SetDefault("server.port", "8090")
	viper.SetDefault("server.grpc_port", "4890")
	viper.SetDefault("server.debug", false)
	viper.SetDefault("server.sentry", false)
	viper.SetDefault("server.session", true)
	viper.SetDefault("server.masterservice", true)
	viper.SetDefault("server.ip", "0.0.0.0")

	// 5) Read config file if available
	if err := viper.ReadInConfig(); err != nil {
		SetLog("config: no settings.toml found, using defaults and ENV")
	} else {
		SetLog("config: loaded " + viper.ConfigFileUsed())
	}

	// 6) Validate essential values
	if err := ValidateConfig(); err != nil {
		return err
	}

	return nil
}

// ValidateConfig performs minimal validation for required parameters.
// It returns an error but never exits the process directly.
func ValidateConfig() error {
	httpPort := strings.TrimSpace(viper.GetString("server.port"))
	grpcPort := strings.TrimSpace(viper.GetString("server.grpc_port"))

	if httpPort == "" {
		return errors.New("config: server.port must not be empty")
	}
	if grpcPort == "" {
		return errors.New("config: server.grpc_port must not be empty")
	}

	if viper.GetBool("server.masterservice") {
		h := strings.TrimSpace(viper.GetString("microservices.masterservice.host"))
		p := strings.TrimSpace(viper.GetString("microservices.masterservice.port"))
		if h == "" || p == "" {
			return errors.New("config: masterservice enabled but host/port missing")
		}
	}
	return nil
}

// EncryptConfigPasswords encrypts plaintext passwords in settings.toml
// using AES-GCM and stores them in the "$2a##<cipher>" format.
// Safe to call repeatedly: already encrypted values are skipped.
func EncryptConfigPasswords() {
	key := GetAesKey()

	dbPwd := viper.GetString("database.password")
	if dbPwd != "" && !strings.HasPrefix(dbPwd, "$2a##") {
		enc, err := EncryptAES(key, dbPwd)
		if err == nil {
			viper.Set("database.password", "$2a##"+enc)
			_ = viper.WriteConfig()
			SetLog("config: database.password encrypted")
		} else {
			SetErrorLog("config: failed to encrypt DB password: " + err.Error())
		}
	}

	emailPwd := viper.GetString("email.password")
	if emailPwd != "" && !strings.HasPrefix(emailPwd, "$2a##") {
		enc, err := EncryptAES(key, emailPwd)
		if err == nil {
			viper.Set("email.password", "$2a##"+enc)
			_ = viper.WriteConfig()
			SetLog("config: email.password encrypted")
		} else {
			SetErrorLog("config: failed to encrypt email password: " + err.Error())
		}
	}
}

// DecryptConfigPasswords decrypts values stored as "$2a##<cipher>".
// Compatible with both legacy (simple AES) and new AES-GCM formats.
func DecryptConfigPasswords(pwd string) string {
	if pwd == "" {
		return ""
	}
	if !strings.HasPrefix(pwd, "$2a##") {
		return pwd
	}

	parts := strings.SplitN(pwd, "##", 2)
	if len(parts) != 2 {
		return pwd
	}

	key := GetAesKey()
	// Try new AES-GCM decryption first
	if dec, err := DecryptAES(key, parts[1]); err == nil {
		return dec
	}

	// Fallback: legacy decrypt() method
	msg, _ := decrypt(parts[1])
	return msg
}

// ConfigString returns a configuration value as string.
func ConfigString(key string) string {
	return viper.GetString(key)
}

// ConfigBool returns a configuration value as bool.
func ConfigBool(key string) bool {
	return viper.GetBool(key)
}

// ConfigInt returns a configuration value as int.
func ConfigInt(key string) int {
	return viper.GetInt(key)
}

// GetPass safely resolves passwords from either ENV or encrypted config values.
// Priority: explicit ENV variable > encrypted TOML value > plaintext fallback.
func GetPass(conf string) string {
	pwd := viper.GetString(conf)

	// 1) Check if there is an ENV override reference (like database.password_env)
	if strings.Contains(conf, "password") {
		envKey := viper.GetString(conf + "_env")
		if envKey != "" {
			if val, ok := os.LookupEnv(envKey); ok && val != "" {
				return val // priority: ENV wins
			}
		}
	}

	// 2) If empty, nothing to decrypt
	if pwd == "" {
		return ""
	}

	// 3) If encrypted, decrypt
	if strings.HasPrefix(pwd, "$2a##") {
		return DecryptConfigPasswords(pwd)
	}

	// 4) Otherwise, return as is (plaintext fallback)
	return pwd
}

// Int32 returns a pointer to int32 (helper for proto structs).
func Int32(v int) *int32 {
	s := int32(v)
	return &s
}

// ConvertInterfaceToAny serializes any Go value into protobuf Any.
func ConvertInterfaceToAny(v interface{}) (*anypb.Any, error) {
	anyValue := &any.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{Value: bytes}
	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	return anyValue, err
}

// ConvertAnyToInterface deserializes a protobuf Any into Go interface{}.
func ConvertAnyToInterface(anyValue *anypb.Any) (interface{}, error) {
	var value interface{}
	bytesValue := &wrappers.BytesValue{}
	err := anypb.UnmarshalTo(anyValue, bytesValue, proto.UnmarshalOptions{})
	if err != nil {
		return value, err
	}
	uErr := json.Unmarshal(bytesValue.Value, &value)
	if uErr != nil {
		return value, uErr
	}
	return value, nil
}

// ToMapStringAny converts map[string]interface{} to map[string]*anypb.Any.
func ToMapStringAny(v map[string]interface{}) map[string]*anypb.Any {
	if len(v) == 0 {
		return nil
	}
	fields := make(map[string]*anypb.Any, len(v))
	for k, val := range v {
		fields[k], _ = ConvertInterfaceToAny(val)
	}
	return fields
}

// ToMapStringInterface converts map[string]*anypb.Any to map[string]interface{}.
func ToMapStringInterface(v map[string]*anypb.Any) map[string]interface{} {
	if len(v) == 0 {
		return nil
	}
	fields := make(map[string]interface{}, len(v))
	for k, val := range v {
		fields[k], _ = ConvertAnyToInterface(val)
	}
	return fields
}
