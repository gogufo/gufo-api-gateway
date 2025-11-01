# 🦉 Gufo API Gateway

**Gufo** (from Italian — *“owl”*) is an open-source, lightweight **gRPC + REST API Gateway**  
for modular microservice architectures.  
Originally designed as a RESTful plugin framework, Gufo has evolved into a secure, flexible, and production-ready gateway  
with a focus on **simplicity**, **security**, and **extensibility**.

---

## ✨ Key Features

- 🧩 **Modular Architecture** — plug in your own gRPC or REST microservices  
- 🔐 **Secure by Default** — AES-GCM encrypted configs, TLS & mTLS support, and environment-based secrets  
- 🚀 **Zero-Config Startup** — Gufo auto-creates a minimal config if missing  
- 🧠 **Master-Service Discovery** — dynamic registration of connected microservices  
- 🪶 **Lightweight Core** — written in pure Go with minimal dependencies  
- 📦 **Docker-Ready** — one-command build and run  
- ⚙️ **Extensible** — easily build your own plugins or sidecar services  

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

````

You can use Gufo in **standalone mode** (no auth/session)  
or as part of a full microservice ecosystem with authentication and permissions.

---

## 🛠️ Quick Start

### 1️⃣ Build the Docker image

```bash
docker build --no-cache -t amyerp/gufo-api-gateway:latest -f Dockerfile .
````

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

Structured JSON logging with daily file rotation and max 10 MB per log file.
Logs are safe to forward to ELK, Loki, or Promtail.

---

## 🧰 Generate gRPC Bindings

To regenerate gRPC stubs for your microservices:

```bash
cd proto
docker run -v $PWD:/defs namely/protoc-all \
  -f microservice.proto -o go/ -l go
```

You can also generate clients for other languages (Ruby, C#, Python, etc.) by changing `-l`.

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

## 🧱 Development Roadmap

* ✅ PR-1: Zero-Config startup, fallback creation, ENV-based secrets
* ✅ PR-2: JSON logging, daily rotation, AES-GCM encryption
* 🔜 PR-3: CLI toolset (key rotation, certificate management)
* 🔜 PR-4: Prometheus / OpenTelemetry metrics
* 🔜 PR-5: REST auto-documentation (Swagger-like)

---

## 💬 Contributing

Gufo is open-source and welcomes contributions!

You can:

* Open issues or pull requests on [GitHub](https://github.com/gogufo/gufo-api-gateway)
* Submit microservice examples
* Improve docs or translations

---

## 📜 License

Licensed under the **Apache License 2.0**
© 2019 – 2025 Alexey Yanchenko. All rights reserved.


