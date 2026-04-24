# 🦉 Gufo API Gateway

[![GoDoc](https://godoc.org/github.com/gogufo/gufo-api-gateway?status.svg)](https://godoc.org/github.com/gogufo/gufo-api-gateway)
[![test status](https://github.com/gogufo/gufo-api-gateway/actions/workflows/build.yml/badge.svg)](https://github.com/gogufo/gufo-api-gateway/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/gogufo/gufo-api-gateway)](https://goreportcard.com/report/github.com/gogufo/gufo-api-gateway)

**Gufo** (from Italian — *“owl”*) is an open-source, lightweight production-grade **gRPC + REST API Gateway**
for modular microservice architectures.
Originally designed as a RESTful plugin framework, Gufo has evolved into a secure, flexible, and production-ready gateway
with a focus on **simplicity**, **security**, and **extensibility**.


Gufo can operate as:
- a **standalone gateway** for small deployments
- or a **cluster gateway** with central MasterService and Redis-backed discovery

Designed for **high-load production environments**, Kubernetes, and secure internal traffic.

Gufo versions <= 1.21 are licensed under the Apache License, Version 2.0.
Gufo versions >= 1.22 are licensed under the Business Source License 1.1 (BSL).

The Change Date for BSL-licensed versions is January 1, 2029.
On that date, the code will automatically be re-licensed under the Apache License, Version 2.0.


---

## ✅ Production Readiness

Gufo is officially verified for production:

📄 **Production Checklist:**  
→ [`/docs/PRODUCTION-CHECKLIST.md`](./docs/PRODUCTION-CHECKLIST.md)

Verified:
- Standalone Mode
- Cluster Mode with MasterService
- Kubernetes
- Rate limiting, metrics, mTLS, CI/CD

---

## ✨ Key Features

* 🧩 **Modular Architecture** — plug in your own gRPC or REST microservices
* 🔐 **Secure by Default** — AES-GCM encrypted configs, TLS & mTLS support, and environment-based secrets
* 🚀 **Zero-Config Startup** — Gufo auto-creates a minimal config if missing
* 🧠 **Master-Service Discovery** — dynamic registration of connected microservices
* 🪶 **Lightweight Core** — written in pure Go with minimal dependencies
* 📦 **Docker-Ready** — one-command build and run
* ⚙️ **Extensible** — easily build your own plugins or sidecar services
* 📡 **gRPC Connection Pool** — TLS/mTLS, retries, deadlines, keepalive, and per-service timeouts
* 📁 **Streaming & Multi-File Upload** — REST `PUT` → gRPC streaming bridge
* 📊 **Metrics & Telemetry** — Prometheus + OpenTelemetry integration

---

## 🧠 Architecture Overview

Gufo acts as a **universal entry point** between REST clients and gRPC microservices.
Each request is validated, optionally authorized, and routed to the appropriate backend module.

```
[ Client ]
↓
[ Gufo Gateway ]
↓
├── Masterservice   — service discovery and microservice registry
├── Auth            — authentication (login/password)
├── Session         — OAuth2.0 token management
├── Rights          — access control, OTP, and API tokens
├── Notifications   — email / chat delivery
├── User            — user profile management
└── Reg             — registration microservice for onboarding
```

You can use Gufo in **standalone mode** (no auth/session)
or as part of a full microservice ecosystem with authentication and permissions.


---

## ⚙️ Operating Modes

### 🟢 Standalone Mode (default, minimal)

Used for:
- Local development
- Edge gateway
- Small deployments

```bash
GUFO_SERVER_MASTERSERVICE=false
````

Characteristics:

* ✅ No Redis
* ✅ No MasterService
* ✅ Heartbeat returns local mock
* ✅ Routing via ENV only
* ✅ Fully production safe

---

### 🔵 Cluster Mode (production microservices)

Used for:

* Kubernetes clusters
* Multi-service platforms

```bash
GUFO_SERVER_MASTERSERVICE=true
```

Requirements:

* ✅ MasterService running
* ✅ Redis mandatory
* ✅ Heartbeat through gRPC
* ✅ Dynamic service registry
* ✅ TTL, sweeper, fallback cache

⚠️ **Redis is required ONLY when `masterservice=true`.**


---

## 🛠️ Quick Start

### 1️⃣ Build the Docker image

```bash
docker build --no-cache -t amyerp/gufo-api-gateway:latest -f Dockerfile .
```

### 2️⃣ Run the container

```bash
docker run --rm -it \
  -e GUFO_SERVER_MASTERSERVICE=false \
  -e GUFO_SECURITY_MODE=sign \
  -e GUFO_SERVER_SIGN=dev \
  -e GUFO_SERVER_PORT=8090 \
  -e GUFO_SERVER_GRPC_PORT=4890 \
  -p 8090:8090 \
  -p 4890:4890 \
  -p 9100:9100 \
  amyerp/gufo-api-gateway:latest
```

If no config file is found, Gufo automatically creates a safe default:

```
config/settings.toml
/etc/gufo/secret.key
```

---

## ✅ Core Production Endpoints

| Endpoint            | Description                                            |
| ------------------- | ------------------------------------------------------ |
| `/api/v1/info`      | **Primary production endpoint** (version, build, mode) |
| `/api/v1/health`    | Kubernetes / Load balancer healthcheck                 |
| `/api/v1/heartbeat` | Microservice registration                              |
| `/metrics`          | Prometheus metrics (token-protected)                   |

### Example

```bash
curl http://127.0.0.1:8090/api/v1/info
```

Response:

```json
{
  "data": {
    "registration": false,
    "version": "1.21.0 (9d0398b, 2025-12-07T04:44:03Z)"
  }
}
```

---
## 📦 Official Docker Image for Gufo API Gateway

Official Gufo API Gateway image is published on Docker Hub:

 https://hub.docker.com/r/amyerp/gufo-api-gateway

Pull the latest version:

```bash
 docker pull amyerp/gufo-api-gateway:latest
```

Tags:

* `latest` — stable
* `vX.Y.Z` — releases
* `dev` — CI auto-build


---
## Kubernetes Deployment

Gufo API Gateway is ready for Kubernetes and exposes dedicated ports for each traffic type:

- **4890** — Internal gRPC (microservices → gateway)
- **8090** — HTTP (health checks / Ingress)
- **9100** — Metrics (Prometheus)

### Apply Deployment and Services

```bash
kubectl apply -f gufo.yml
```
---
### 🧰 Manual Installation (without Docker)

You can run Gufo directly from source without using Docker.

### 🧩 Requirements

- **Go 1.25+**
- **Redis** — optional, used for session caching
- **OpenSSL** — required only if mTLS is enabled
- **Linux or macOS** environment recommended

### ⚙️ Installation

```bash
git clone https://github.com/gogufo/gufo-api-gateway.git
cd gufo-api-gateway
go build -o gufo gufo.go
sudo ./gufo start
```

### Quick Test

```go
mkdir -p /var/gufo/config
cp config/settings.example.toml /var/gufo/config/settings.toml
./gufo

```
---

## ⚙️ Configuration

Gufo supports **layered configuration**:

| Priority | Source                    | Description                                            |
| -------- | ------------------------- | ------------------------------------------------------ |
| 1️⃣      | **Environment variables** | Highest priority, ideal for Docker and CI/CD           |
| 2️⃣      | **.env file**             | Optional local development overrides                   |
| 3️⃣      | **settings.toml**         | Default configuration file (auto-generated if missing) |

### Example environment variables

```bash
GUFO_DB_PASS=supersecret
GUFO_SIGN=my_internal_sign_key
GUFO_AES_KEY=my_encryption_key
GUFO_SENTRY_DSN=https://<your_sentry_dsn>
```

### Default `settings.toml`

```toml
[server]
port = "8090"
grpc_port = "4890"
debug = false
sentry = false
session = true
masterservice = true
ip = "0.0.0.0"
sysdir = "/var/gufo/"
tempdir = "/var/gufo/templates/"
filedir = "/var/gufo/files/"
plugindir = "/var/gufo/lib/"
logdir = "/var/gufo/log/"

[database]
type = "mysql"
host = "db"
port = "3306"
dbname = "gufo"
user = "root"
password_env = "GUFO_DB_PASS"

[redis]
host = "redis://redis"

[microservices.masterservice]
host = "masterservice"
port = "5300"
type = "server"
entrypointversion = "1.0.0"
cron = false
```
---

## 🧩 Generate gRPC Connection Files

> _Tip: Always re-generate gRPC bindings after updating `microservice.proto`._


Before building or integrating new microservices, regenerate the Go bindings from the `.proto` schema.

Go to the `/proto` folder and run:

```bash
docker run -v $PWD:/defs namely/protoc-all \
  -f microservice.proto -o go/ -l go
```
You can replace -l go with another language:
ruby, python, csharp, node, php, etc.

All generated gRPC files will be placed in /proto/go/

---

## Control Plane & Dual-Mode Routing

The Gateway supports a safe **dual-mode operation**:

- **Cluster Mode** (`server.masterservice = true`):
    - All routing, heartbeat and cron control are proxied via **MasterService**.
    - Gateway dynamically resolves microservice hosts through MasterService.
    - Timeout and fallback to local registry are applied for fault tolerance.

- **Standalone Mode** (`server.masterservice = false`):
    - Gateway resolves microservice hosts directly from environment configuration.
    - MasterService is fully optional.
    - Heartbeat requests return a local mock response (`leader=true`, `cron=true`).

### Reliability Improvements
- Eliminated unsafe request mutation during MasterService resolution.
- Added timeout protection for MasterService calls.
- Removed registry as a single point of failure by adding dynamic fallback resolution.

This design keeps microservices **mode-agnostic**, moves all control logic into the Gateway, and fully supports both small standalone deployments and clustered production environments without code changes in microservices.


## 🔐 Security Model

Gufo implements several layers of protection:

### 1️⃣ AES-GCM Configuration Encryption

All sensitive data (passwords, tokens, secrets) are automatically encrypted with AES-GCM.
Each installation stores its unique key at `/etc/gufo/secret.key` or via `GUFO_AES_KEY`.

### 2️⃣ mTLS and Internal Auth Signatures

Internal gRPC communication can be authenticated by:

* a shared system sign (`GUFO_SIGN`),
* or full **mutual TLS** (mTLS) between the gateway and microservices.

### 3️⃣ Error Isolation

Each service runs independently — gateway failures never expose credentials or plaintext configs.

### 4️⃣ Logging and Rotation

Structured JSON logging with daily rotation and safe forwarding to ELK/Loki/Promtail.

---

## 📡 REST → gRPC Streaming & File Uploads

Gufo automatically converts REST `PUT` requests into gRPC streaming calls.
This allows file uploads (single or multiple) to be **directly streamed** from the REST client to a backend gRPC microservice —
without ever storing files on the Gateway itself.

### Example: Binary Upload

```bash
curl -X PUT http://localhost:8090/api/v3/storage/upload \
  -H "X-Filename: demo.jpg" \
  -H "Content-Type: application/octet-stream" \
  --data-binary "@demo.jpg"
```

The Gateway:

1. Detects the `PUT` method
2. Sets `t.IR.Param = "stream"`
3. Calls `GRPCStreamPut()`
4. Streams the body in 64 KB chunks via gRPC

---

### Example: Multi-File Upload

```bash
curl -X PUT http://localhost:8090/api/v3/storage/upload \
  -F "file=@file1.jpg" \
  -F "file=@file2.jpg"
```

Each file is streamed independently to the destination gRPC service.

---

### Example: Server-Side Handler

```go
func (s *Server) Stream(stream pb.Reverse_StreamServer) error {
	var currentFile string
	var f io.WriteCloser

	for {
		req, err := stream.Recv()
		if err == io.EOF { return nil }
		if err != nil { return err }

		if anyChunk, ok := req.Args["chunk"]; ok {
			var chunk pb.FileChunk
			if err := anypb.UnmarshalTo(anyChunk, &chunk, proto.UnmarshalOptions{}); err != nil {
				return err
			}

			if f == nil {
				currentFile = chunk.Name
				f, err = os.Create(filepath.Join("/tmp", currentFile))
				if err != nil { return err }
			}

			if len(chunk.Data) == 0 {
				f.Close()
				f = nil
				continue
			}

			f.Write(chunk.Data)
		}
	}
}
```

---

## 🔄 gRPC Connection Pool

Located in `gufodao/grpcpool.go`, the connection pool provides:

* Persistent `sync.Map` of `host:port → *grpc.ClientConn`
* TTL: 5 minutes
* Background sweeper that closes expired connections
* TLS / mTLS support
* Retry policy (`WithDefaultServiceConfig`) — up to 4 attempts on `UNAVAILABLE`
* Keepalive every 30 seconds
* Per-service timeouts (via `microservices.<name>.timeout`)

```toml
[microservices.storage]
host = "127.0.0.1"
port = "4802"
timeout = "8s"
stream_timeout = "2m"
```

---

## 🔌 Transport Abstraction

Gufo supports pluggable transports via `transport.Transport` interface.  
The default implementation is `GRPCTransport`, but you can register your own:

```go
transport.Register(&MyCustomTransport{})
```

---

## 🧩 CLI Commands

| Command               | Description                                    |
| --------------------- | ---------------------------------------------- |
| `gufo start`          | Start API Gateway                              |
| `gufo stop`           | Stop running instance                          |
| `gufo cert init`      | Generate self-signed TLS certificates          |
| `gufo key rotate`     | Rotate encryption key and re-encrypt passwords |
| `gufo migrate config` | Migrate legacy config to new AES-GCM format    |

---


## 🔄 Service Registry & Heartbeat

```bash
POST /api/v1/heartbeat
{
  "service": "auth",
  "host": "auth",
  "port": "5301"
}
```

Cluster mode:

* forwarded to MasterService
* cached with TTL
* fallback used on failure

Standalone:

* local mock response

---

## 🧱 Middleware Chain

* Request ID
* Logger
* CORS
* Rate Limiter (RPS + Burst)

Returns:

* `HTTP 429` on limit exceed

---

## 📊 Metrics & Observability

Gufo includes built-in **Prometheus** and **OpenTelemetry** instrumentation for runtime visibility.

### Prometheus Metrics

Metrics are exposed on port `9100`:

```

[http://127.0.0.1:9100/api/v3/metrics](http://127.0.0.1:9100/api/v3/metrics)

````

Access is protected via a token:

```bash
curl -H "X-Metrics-Token: gufo-metrics" http://127.0.0.1:9100/api/v3/metrics
````

### Available Metrics

| Metric                               | Description                       |
| ------------------------------------ | --------------------------------- |
| `gufo_http_requests_total`           | Total number of HTTP requests     |
| `gufo_http_request_duration_seconds` | Histogram of request latency      |
| `gufo_grpc_pool_hits_total`          | gRPC connection pool cache hits   |
| `gufo_grpc_pool_misses_total`        | gRPC connection pool cache misses |
| `gufo_grpc_retries_total`            | Number of gRPC retry attempts     |

### OpenTelemetry Tracing

To enable distributed tracing, set in your config:

```toml
[server]
telemetry = true
```

Each request automatically propagates `trace-id` and `span-id`
through HTTP → gRPC → microservices, enabling full end-to-end observability.



---

## 🔄 Service Registry and Heartbeat

Gufo introduces a lightweight **in-memory service registry** that keeps track of all registered microservices  
and ensures continuous routing — even if the central **MasterService** becomes temporarily unavailable.

---

### 🧠 How It Works

| Step | Description |
|------|--------------|
| 1️⃣ | **Initial discovery** — when a module is requested for the first time, Gufo queries `masterservice` via gRPC (`getmicroservicebypath`) to get its host and port. |
| 2️⃣ | **Registry cache** — the response is stored locally (`registry.ServiceInfo`) with a time-to-live (TTL). A background refresher periodically revalidates each entry. |
| 3️⃣ | **Failover** — if `masterservice` is unavailable, Gufo automatically uses the last valid cached entry to continue serving requests. |
| 4️⃣ | **Heartbeat endpoint** — microservices periodically send heartbeats to inform the MasterService that they are alive and healthy. |
| 5️⃣ | **Self-healing** — inactive or non-reporting services are automatically purged from the registry after TTL expiration. |

---

### ⚡ Heartbeat Endpoint

Microservices should periodically send a JSON heartbeat to Gufo’s REST endpoint:

```bash
POST /api/v3/heartbeat
Content-Type: application/json

{
  "service": "auth",
  "host": "auth",
  "port": "5301"
}
```

Gufo will automatically forward this heartbeat to the masterservice through gRPC,
updating the registry and refreshing the LastUpdate timestamp.

```toml
[masterregistry]
ttl = "60s"              # time-to-live for cached entries
refresh_interval = "30s" # background refresh interval
sweeper_interval = "1m"  # periodic cleanup of expired entries
```

✅ Benefits
* 🔁 Continuous operation even if masterservice is offline
* ⚡ Low-latency service lookup via local cache
* 🧩 Seamless integration with existing gRPC routing
* 🧠 Intelligent self-cleaning and background refresh
* 🔐 All communication still flows securely through Gufo’s /api/v3/* endpoints

---

## 🧩 Related Microservices and Tools

Gufo works best as part of the **microservice ecosystem**, where each component  
has a clearly defined role and communicates through secure gRPC channels.

| Service | Repository | Description |
|----------|-------------|-------------|
| 🧭 **MasterService** | [gogufo/masterservice](https://github.com/gogufo/masterservice) | Central service registry and discovery endpoint. Handles microservice registration, heartbeats, and health monitoring. |
| 🔐 **Session Service** | [gogufo/session-m10e](https://github.com/gogufo/session-m10e) | Handles user sessions, tokens, and permissions management for all connected clients. |
| ⚙️ **Microservice Generator** | [gogufo/gufo-grpc-microservice-generator](https://github.com/gogufo/gufo-grpc-microservice-generator) | CLI + Docker tool to scaffold new Gufo-compatible microservices in seconds. |

---

### 🚀 Generate a New Microservice (via Docker)

You can instantly create a new Gufo-compatible microservice skeleton  
using the **official generator image**:

```bash
docker run --rm -v $(pwd):/src amyerp/gufo_grpc_microservice_generator:latest /go/bin/grpccreator my_new_service
```

This command will:

* scaffold a ready-to-build Go microservice project under the current directory

* include boilerplate gRPC definitions, Dockerfile, Makefile, and config templates

* automatically register it with Gufo Gateway (via /api/v3/heartbeat) once running

---

## 🧱 Middleware Framework

Gufo includes a lightweight middleware system that processes every request  
**before** and **after** it reaches the core gateway.

Built-in middlewares:
- 🪶 `RequestID` — adds `X-Request-ID` to each request
- 🪵 `Logger` — structured logging (method, path, latency, status)
- 🌍 `CORS` — standard CORS headers for browser APIs
- ⚖️ `RateLimiter` — simple in-memory token bucket

Usage:
```go
middleware.Register(middleware.NewRequestID())
middleware.Register(middleware.NewLogger())
middleware.Register(middleware.NewCORS())
middleware.Register(middleware.NewRateLimiter(100, time.Second))
```

Each middleware can modify context, block requests, or log responses.
Custom middleware can be added by implementing the Middleware interface.


## 💬 Contributing

Gufo is open-source and welcomes contributions!

You can:

* Open issues or pull requests on [GitHub](https://github.com/gogufo/gufo-api-gateway)
* Submit microservice examples
* Improve documentation or translations

---

## 📜 License

Licensed under the **Business Source License 1.1**
© 2019–2026 Alexey Yanchenko. All rights reserved.

