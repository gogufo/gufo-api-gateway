// Copyright 2025 Alexey Yanchenko
//
// Graceful shutdown helper for Gufo.
// Catches OS signals and executes cleanup callbacks safely.

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
