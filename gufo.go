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
//
//
// This is main file, from which starts Gufo

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	"github.com/gogufo/gufo-api-gateway/registry"
	"github.com/gogufo/gufo-api-gateway/transport"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/certifi/gocertifi"
	handler "github.com/gogufo/gufo-api-gateway/handler"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	v "github.com/gogufo/gufo-api-gateway/version"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	viper "github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()
var grpcSrv *grpc.Server

// info is function for CLI.
// in this function determinate to start Web server
func info() {

	app.Name = "Gufo API Gateway"
	app.Usage = "RESTful API with GRPC microservices"
	app.Version = v.VERSION
	app.Action = StartService

}

// commands is function for CLI
func commands() {
	app.Commands = []*cli.Command{
		{
			Name:   "stop",
			Usage:  "Stop Gufo Server",
			Action: StopApp,
		},
		{
			Name:  "cert",
			Usage: "Certificate management commands",
			Subcommands: []*cli.Command{
				{
					Name:   "init",
					Usage:  "Generate self-signed CA, server, and client certificates for mTLS",
					Action: sf.GenerateCertificates,
				},
			},
		},
	}
}

// main is the entry point of Gufo API Gateway.
// It initializes the configuration, sets up Sentry (if enabled),
// prepares CLI commands, and starts the main service.

func main() {
	// Initialize configuration
	sf.EnsureConfigExists()
	if err := sf.InitConfig(); err != nil {
		sf.SetErrorLog("config init failed: " + err.Error())
		os.Exit(1)
	}
	sf.EncryptConfigPasswords()

	ctx := context.Background()
	sf.InitTelemetry(ctx)
	defer sf.ShutdownTelemetry(ctx)

	// Initialize Sentry (optional)
	if viper.GetBool("server.sentry") {
		sf.SetLog("Connecting to Sentry...")

		sentryClientOptions := sentry.ClientOptions{
			Dsn:              viper.GetString("sentry.dsn"),
			EnableTracing:    viper.GetBool("sentry.tracing"),
			Debug:            viper.GetBool("sentry.debug"),
			TracesSampleRate: viper.GetFloat64("sentry.trace"),
		}

		// Load trusted CA certificates
		rootCAs, err := gocertifi.CACerts()
		if err != nil {
			sf.SetLog("Could not load CA certificates for Sentry: " + err.Error())
		} else {
			sentryClientOptions.CaCerts = rootCAs
		}

		// Initialize Sentry client
		err = sentry.Init(sentryClientOptions)
		if err != nil {
			sf.SetLog("Error initializing Sentry: " + err.Error())
		} else {
			flushsec := viper.GetDuration("sentry.flush")
			defer sentry.Flush(flushsec * time.Second)
		}
	}

	// Setup CLI metadata and commands
	info()
	commands()

	// Run CLI and web server
	err := app.Run(os.Args)
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			sf.SetErrorLog("gufo.go: main: " + err.Error())
		}
	}
}

// StopApp gracefully stops a running Gufo instance via the /exit endpoint.
// This ensures that telemetry, Sentry, and all servers shut down cleanly.
// Works both in Docker and bare-metal setups.
func StopApp(c *cli.Context) error {
	addr := viper.GetString("server.ip")
	port := viper.GetString("server.port")

	if addr == "" {
		addr = "127.0.0.1"
	}
	if port == "" {
		port = "8090"
	}

	url := fmt.Sprintf("http://%s:%s/exit", addr, port)
	sf.SetLog("CLI command 'gufo stop' → sending shutdown signal to " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		sf.SetErrorLog("stop: cannot create request: " + err.Error())
		return err
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		sf.SetErrorLog("stop: cannot reach Gufo server: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		sf.SetErrorLog(fmt.Sprintf("stop: server responded with %s", resp.Status))
		fmt.Printf("Stop command failed: %s\n", resp.Status)
		return fmt.Errorf("server responded with %s", resp.Status)
	}

	fmt.Println("✅ Shutdown signal sent — Gufo is stopping gracefully...")
	sf.SetLog("Shutdown request acknowledged by Gufo")
	return nil
}

// ExitApp handles graceful stop via HTTP request (debug mode only).
func ExitApp(w http.ResponseWriter, r *http.Request) {
	if !viper.GetBool("server.debug") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	sf.SetLog("Received /exit request — initiating graceful shutdown")

	go func() {
		time.Sleep(500 * time.Millisecond) // small delay to let HTTP 200 return
		sf.WaitForShutdown(func() {
			// This callback will stop servers gracefully (same as SIGTERM)
			sf.SetLog("Graceful shutdown triggered by /exit endpoint")
		})
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Gufo is shutting down..."))
}

// StartService is function for start WEB Server to listen port
func StartService(c *cli.Context) (rtnerr error) {
	// Initialize Redis
	sf.InitCache()

	registry.StartRefresher()
	sf.SetLog("🧠 Registry cache refresher started")

	// Register default transport (gRPC)
	transport.Register(&transport.GRPCTransport{})
	sf.SetLog("✅ Registered default transport: gRPC")

	port := sf.ConfigString("server.port")

	m := fmt.Sprintf("Gufo v%s starting on :%s (gRPC :%s, mode=%s)",
		v.VERSION,
		viper.GetString("server.port"),
		viper.GetString("server.grpc_port"),
		strings.ToLower(viper.GetString("security.mode")),
	)

	sf.SetLog(m)
	fmt.Printf(m)

	// ---------------------------------------------------
	// Router initialization (replaces http.Handle*)
	// ---------------------------------------------------
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(sf.RecoveryMiddleware)          // panic-safe middleware
	r.Use(otelhttp.NewMiddleware("gufo")) // telemetry tracing

	// Routes
	r.Route("/api/v3", func(r chi.Router) {
		r.Get("/health", handler.Health)
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.API(w, r, 3)
		}))
	})

	if viper.GetBool("server.debug") {
		r.Get("/exit", ExitApp)
	}

	// ---------------------------------------------------
	// Start servers
	// ---------------------------------------------------
	go StartGRPCService()

	// Internal metrics server (localhost only, protected by X-Metrics-Token)
	go func() {
		metricsMux := http.NewServeMux()
		token := viper.GetString("server.metrics_token")

		metricsMux.HandleFunc("/api/v3/metrics", func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				http.Error(w, "Metrics endpoint disabled", http.StatusForbidden)
				return
			}
			if r.Header.Get("X-Metrics-Token") != token {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			promhttp.Handler().ServeHTTP(w, r)
		})

		if err := http.ListenAndServe(":9100", metricsMux); err != nil {
			sf.SetErrorLog("metrics server error: " + err.Error())
		}
	}()

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	sf.SetLog(fmt.Sprintf("HTTP listening on :%s", port))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sf.SetErrorLog("HTTP server error: " + err.Error())
		}
	}()

	sf.WaitForShutdown(func() {
		if grpcSrv != nil {
			grpcSrv.GracefulStop()
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	})

	return nil
}

func StartGRPCService() {
	getport := strings.TrimSpace(viper.GetString("server.grpc_port"))
	port := ":4890"
	if getport != "" {
		port = ":" + getport
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	sf.SetLog(fmt.Sprintf("gRPC listening on %s", port))

	var opts []grpc.ServerOption

	if viper.GetBool("server.grpc_tls_enabled") {
		certPath := viper.GetString("security.cert_path")
		keyPath := viper.GetString("security.key_path")
		caPath := viper.GetString("security.ca_path")

		serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			grpclog.Fatalf("cannot load server cert: %v", err)
		}

		caBytes, err := os.ReadFile(caPath)
		if err != nil {
			grpclog.Fatalf("cannot read CA file: %v", err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caBytes)

		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{serverCert},
			ClientCAs:    caPool,
		}
		if viper.GetBool("server.grpc_mtls_enabled") {
			tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
			sf.SetLog("gRPC mTLS mode enabled")
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsCfg)))
	}

	grpcSrv = grpc.NewServer(opts...)
	pb.RegisterReverseServer(grpcSrv, &Server{})

	if err := grpcSrv.Serve(listener); err != nil {
		sf.SetErrorLog("gRPC server error: " + err.Error())
	}
}

type Server struct {
}

// Do handles incoming gRPC requests and verifies authentication
func (s *Server) Do(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {
	case "hmac":
		secret := viper.GetString("security.hmac_secret")
		maxAge := time.Duration(viper.GetInt("security.max_age")) * time.Second

		// Extra safety check: avoid nil dereference in request.Module or request.Sign
		if request.Sign == nil || request.Module == nil ||
			!sf.VerifyHMAC(secret, *request.Module, *request.Sign, maxAge) {

			sf.SetErrorLog("Unauthorized gRPC request (HMAC mode)")
			return sf.ErrorReturn(request, 401, "00001", "Invalid or expired signature"), nil
		}

	case "sign":
		if request.Sign == nil || viper.GetString("server.sign") != *request.Sign {
			sf.SetErrorLog("Unauthorized gRPC request (static sign mode)")
			return sf.ErrorReturn(request, 401, "00001", "Invalid signature"), nil
		}

	case "mtls":
		// For mTLS mode, trust gRPC layer verification — no Sign check needed
		sf.SetLog("mTLS mode active - skipping sign verification")

	default:
		sf.SetErrorLog("Unknown security mode")
		return sf.ErrorReturn(request, 500, "00002", "Security mode not configured"), nil
	}

	return handler.InternalRequest(request), nil
}

// Stream handles bidirectional streaming RPC calls.
func (s *Server) Stream(stream pb.Reverse_StreamServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		val, _ := anypb.New(wrapperspb.String("pong"))

		resp := &pb.Response{
			Data: map[string]*anypb.Any{
				"echo": val,
			},
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
