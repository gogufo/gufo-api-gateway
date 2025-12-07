# Gufo mTLS Initialization in Kubernetes (k3s)

This document describes the **correct and production-safe procedure** for initializing and distributing **mTLS certificates** for **Gufo API Gateway** and all connected microservices in a Kubernetes (k3s) cluster.

---

## 1. What `gufo cert init` Generates

The command:

```bash
gufo cert init
```

generates a full **mTLS trust chain**:

| File             | Purpose                              |
| ---------------- | ------------------------------------ |
| `ca.pem`         | Root Certificate Authority           |
| `server.pem`     | Gufo server certificate              |
| `server-key.pem` | Gufo server private key              |
| `client.pem`     | Client certificate for microservices |
| `client-key.pem` | Client private key                   |

This is **not a simple self-signed cert** — this is a full **CA → Server → Client** chain.

---

## 2. Why You Must NOT Generate Certificates Inside Kubernetes Pods

Running `gufo cert init` **inside a Pod** is **architecturally incorrect**:

* ❌ Certificates exist only in the container filesystem
* ❌ New Pods (HPA, restarts) will NOT have the same certs
* ❌ Microservices will lose trust on Pod recreation
* ❌ mTLS will silently break in production
* ❌ Cluster becomes non-deterministic

➡️ **Certificates must be generated ONCE outside the cluster and stored as Kubernetes Secrets.**

---

## 3. Correct Production Flow (Recommended)

### ✅ Step 1. Generate Certificates Once (Outside the Cluster)

On your local machine or in CI/CD:

```bash
docker run --rm -v $(pwd)/certs:/certs \
  amyerp/gufo-api-gateway:latest \
  gufo cert init
```

Result:

```
certs/
  ca.pem
  server.pem
  server-key.pem
  client.pem
  client-key.pem
```

---

### ✅ Step 2. Create Kubernetes Secret

```bash
kubectl create secret generic gufo-mtls \
  --from-file=ca.pem \
  --from-file=server.pem \
  --from-file=server-key.pem \
  --from-file=client.pem \
  --from-file=client-key.pem
```

Verify:

```bash
kubectl get secret gufo-mtls -o yaml
```

---

### ✅ Step 3. Mount Certificates into Gufo Pod

In `gufo.yml`:

```yaml
volumes:
  - name: gufo-certs
    secret:
      secretName: gufo-mtls
```

```yaml
volumeMounts:
  - name: gufo-certs
    mountPath: /etc/gufo/certs
    readOnly: true
```

Environment variables for Gufo:

```yaml
env:
  - name: GUFO_SECURITY_MODE
    value: mtls
  - name: GUFO_SECURITY_CERT_PATH
    value: /etc/gufo/certs/server.pem
  - name: GUFO_SECURITY_KEY_PATH
    value: /etc/gufo/certs/server-key.pem
  - name: GUFO_SECURITY_CA_PATH
    value: /etc/gufo/certs/ca.pem
```

---

## 4. How to Provide Certificates to Microservices

### ✅ Recommended: Use the Same Kubernetes Secret

In microservice Deployment:

```yaml
volumes:
  - name: gufo-client-certs
    secret:
      secretName: gufo-mtls
```

```yaml
volumeMounts:
  - name: gufo-client-certs
    mountPath: /etc/gufo/certs
    readOnly: true
```

Microservice environment variables:

```yaml
env:
  - name: GRPC_CLIENT_CERT
    value: /etc/gufo/certs/client.pem
  - name: GRPC_CLIENT_KEY
    value: /etc/gufo/certs/client-key.pem
  - name: GRPC_CLIENT_CA
    value: /etc/gufo/certs/ca.pem
```

This allows the microservice to authenticate itself via **mTLS** when connecting to Gufo.

---

## 5. Full Certificate Lifecycle in Production

```
[ CI / Local Dev Machine ]
          |
          |  gufo cert init
          v
      [ Certificate Files ]
          |
          | kubectl create secret
          v
[ Kubernetes Secret: gufo-mtls ]
          |
          +-------------------------+
          |                         |
   [ Gufo Pod ]               [ Microservice Pods ]
   server.pem                 client.pem
   server-key.pem             client-key.pem
   ca.pem                     ca.pem
```

---

## 6. Security Notes

* ✅ `client-key.pem` and `server-key.pem` **MUST be stored only in Kubernetes Secrets**
* ❌ **Never store private keys in ConfigMaps**
* ✅ `ca.pem` may be stored in ConfigMap if required (read-only trust anchor)
* ✅ All mounts must be read-only

---

## 7. What Happens If You Ignore This Procedure

If certificates are generated dynamically inside Pods:

* ❌ All existing microservices will immediately lose trust
* ❌ Horizontal Pod Autoscaler will create broken Gufo replicas
* ❌ Rolling updates will cause TLS desync
* ❌ Production outages become inevitable

---

## 8. Minimal README Snippet (Optional)

```md
### mTLS in Kubernetes

1. Generate certificates once:
   docker run --rm -v ./certs:/certs amyerp/gufo-api-gateway gufo cert init

2. Create Kubernetes secret:
   kubectl create secret generic gufo-mtls --from-file=certs/

3. Mount secret into Gufo and microservices:
   - /etc/gufo/certs/server.pem
   - /etc/gufo/certs/server-key.pem
   - /etc/gufo/certs/client.pem
   - /etc/gufo/certs/client-key.pem
   - /etc/gufo/certs/ca.pem
```
