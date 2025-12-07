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
