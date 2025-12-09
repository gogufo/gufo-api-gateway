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
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WaitForShutdown listens for SIGINT/SIGTERM and runs the provided callback.
// It gives graceful shutdown to all servers (HTTP, gRPC, metrics, etc.).
func WaitForShutdown(onShutdown func()) {
	// Create channel to listen for system signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	sig := <-sigCh
	SetLog("Received signal: " + sig.String() + " — starting graceful shutdown")

	// Run shutdown callback with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		onShutdown()
		close(done)
	}()

	select {
	case <-done:
		SetLog("Shutdown complete")
	case <-ctx.Done():
		SetErrorLog("Shutdown timed out — forcing exit")
	}

	os.Exit(0)
}
