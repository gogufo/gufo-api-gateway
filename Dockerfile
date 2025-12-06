# ======================
# Stage 1 — Build
# ======================
FROM golang:1.25 AS builder

# Install build deps for CGO and TLS
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential clang git ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

# Enable CGO (needed for TLS)
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=clang

# Build optimized binary
RUN go build -trimpath -ldflags="-s -w" -o /go/bin/gufo gufo.go

# ======================
# Stage 2 — Runtime
# ======================
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy runtime assets and binary
COPY --from=builder /go/bin/gufo /usr/local/bin/gufo
COPY --from=builder /app/config/settings.example.toml /var/gufo/config/settings.toml
COPY --from=builder /app/var/ /var/gufo/

WORKDIR /var/gufo/

EXPOSE 8090 4890 9100

ENTRYPOINT ["gufo"]
