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
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	viper "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// --- Sign / HMAC verification ---

// ComputeHMAC generates an HMAC-based signature using secret, module, and timestamp.
func ComputeHMAC(secret, module string, ts int64) string {
	data := fmt.Sprintf("%s:%d", module, ts)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMAC validates the HMAC signature.
func VerifyHMAC(secret, module, sign string, maxAge time.Duration) bool {
	parts := strings.Split(sign, ":")
	if len(parts) != 2 {
		return false
	}
	sig := parts[0]
	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	if time.Since(time.Unix(ts, 0)) > maxAge {
		return false
	}
	expected := ComputeHMAC(secret, module, ts)
	return hmac.Equal([]byte(sig), []byte(expected))
}

// --- TLS credentials (for mTLS mode) ---

// LoadMTLSCredentials loads client-side TLS credentials for gRPC connections.
func LoadMTLSCredentials() (credentials.TransportCredentials, error) {
	caCertPath := viper.GetString("security.ca_cert")
	clientCertPath := viper.GetString("security.cert")
	clientKeyPath := viper.GetString("security.key")

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append CA certs")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// --- Utility: GetGRPCCredentials chooses credentials based on mode ---

// GetGRPCCredentials returns proper transport credentials depending on security mode.
/*
func GetGRPCCredentials() (grpc.DialOption, error) {
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {
	case "mtls":
		creds, err := LoadMTLSCredentials()
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(creds), nil

	default:
		// for sign and hmac modes, use insecure (plain gRPC)
		return grpc.WithInsecure(), nil
	}
}
*/
func GetGRPCCredentials() (grpc.DialOption, error) {
	mode := strings.ToLower(viper.GetString("security.mode"))

	if mode == "mtls" {
		// mTLS branch
		//return grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)), nil
		creds, err := LoadMTLSCredentials()
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(creds), nil
	}
	
	return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
}
