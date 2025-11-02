# 🦉 Gufo API Gateway

[![GoDoc](https://godoc.org/github.com/gogufo/gufo-api-gateway?status.svg)](https://godoc.org/github.com/gogufo/gufo-api-gateway)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/gogufo/gufo-api-gateway)](https://goreportcard.com/report/github.com/gogufo/gufo-api-gateway)

**Gufo** (from Italian — *“owl”*) is an open-source, lightweight **gRPC + REST API Gateway**
for modular microservice architectures.
Originally designed as a RESTful plugin framework, Gufo has evolved into a secure, flexible, and production-ready gateway
with a focus on **simplicity**, **security**, and **extensibility**.

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

## 🛠️ Quick Start

### 1️⃣ Build the Docker image

```bash
docker build --no-cache -t amyerp/gufo-api-gateway:latest -f Dockerfile .
```

### 2️⃣ Run the container

```bash
docker run -p 8090:8090 \
  -v $(pwd)/config:/config \
  amyerp/gufo-api-gateway:latest
```

If no config file is found, Gufo automatically creates a safe default:

```
config/settings.toml
/etc/gufo/secret.key
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

## 📊 Metrics & Observability

Gufo exposes Prometheus metrics on port `9100`:

```
http://127.0.0.1:9100/api/v3/metrics
```

Access is protected via:

```bash
curl -H "X-Metrics-Token: gufo-metrics" http://127.0.0.1:9100/api/v3/metrics
```

Additionally, OpenTelemetry tracing can be enabled via config (`server.telemetry = true`).

---

## 🧱 Development Roadmap

* ✅ PR-1: Zero-Config startup, fallback creation, ENV-based secrets
* ✅ PR-2: JSON logging, daily rotation, AES-GCM encryption
* ✅ PR-3: CLI toolset (key rotation, certificate management)
* ✅ **PR-4: gRPC Connection Pool, TLS/mTLS, Streaming, File Uploads, Timeouts**
* 🔜 PR-5: REST auto-documentation (Swagger/OpenAPI)

---

## 💬 Contributing

Gufo is open-source and welcomes contributions!

You can:

* Open issues or pull requests on [GitHub](https://github.com/gogufo/gufo-api-gateway)
* Submit microservice examples
* Improve documentation or translations

---

## 📜 License

Licensed under the **Apache License 2.0**
© 2019–2025 Alexey Yanchenko. All rights reserved.

