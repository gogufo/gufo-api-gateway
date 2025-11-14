# **Gufo Platform — White Book (Architecture & Principles)**

### *Version 1.0 — 2025*

**Author:** Alexey Yanchenko

**Status:** Public / Technical Audience

---

# **Table of Contents**

1. Overview
2. Mission and Philosophy
3. Architectural Principles
4. High-Level Architecture
5. Core Components

    * 5.1 API Gateway
    * 5.2 MasterService
    * 5.3 Session Service
    * 5.4 Auth Service
    * 5.5 Registration Service (Reg)
    * 5.6 Notifications Service
    * 5.7 Rights Service
6. Internal Communication Model
7. Security Architecture
8. Data Architecture
9. Geo-Distribution Model
10. Scaling Strategy
11. Deployment Modes
12. Extensibility Model
13.Glossary

---

# 1. **Overview**

**Gufo** is a modular, secure, and geo-distributable microservices platform designed for high-load enterprise applications.
It provides a unified internal protocol, strong security, and clearly separated responsibilities between core services.

Gufo solves the problem of building **scalable, customizable business systems** by providing a ready-made foundation:

* Identity and access control
* Sessions
* Notifications
* User registration
* Microservices registry
* API gateway
* Internal communication security
* Job coordination

Gufo allows businesses to focus on domain logic instead of infrastructure.

---

# 2. **Mission and Philosophy**

Gufo follows several guiding principles:

### 🟦 Separation of Responsibilities

Each service is isolated around **one clear domain**.
No service mixes logic belonging to different business areas.

### 🟩 Horizontal Scalability First

Gufo is designed to scale by **replicating microservices**, not by “making one service bigger”.

### 🟧 Geo-Distributed Architecture

The platform supports local per-region deployments:

* local sessions
* local caching
* local notification providers
* local failover

### 🟥 Replaceability

Each Gufo component can be:

* replaced,
* overridden,
* extended,
* disabled (if the customer does not need it).

### 🟫 Zero Shared Database Requirement

Each service owns **its own database or storage**.
There is no global monolithic DB.

---

# 3. **Architectural Principles**

1. **Control Plane vs Data Plane separation**

    * MasterService = control plane
    * Identity/Session/Auth = data plane

2. **Security-first internal communication**
   Unified protocol with:

    * HMAC
    * mTLS
    * Static Sign token

3. **Locality Matters**
   Sessions are always local to region → low latency.

4. **Customizable Business Flows**
   Registration and notifications vary dramatically across industries → must be separate.

5. **Replaceable Features**
   If a client doesn’t need `reg`, `rights`, or complex notifications — they can be disabled.

---

# 4. **High-Level Architecture**

```
                         ┌───────────────────────┐
                         │     API Gateway       │
                         │ (REST/WS → gRPC, SEC) │
                         └──────────┬────────────┘
                                    │
                              Internal gRPC
                                    │
                         ┌──────────▼────────────┐
                         │     MasterService     │
                         │  (Registry, Routing)  │
                         └───┬─────────┬─────────┘
             ┌────────────────┘         └────────────────┐
             ▼                                          ▼
     ┌──────────────┐                         ┌─────────────────┐
     │   Session     │                         │     Auth        │
     │ (Redis/local) │                         │  (CPU heavy)    │
     └──────────────┘                         └─────────────────┘
             ▼                                          ▼
     ┌──────────────┐                         ┌─────────────────┐
     │ Registration  │                         │    Rights       │
     │ (Custom flow) │                         │ (RBAC/ABAC/etc) │
     └──────────────┘                         └─────────────────┘
                                    ▼
                         ┌─────────────────────┐
                         │   Notifications     │
                         │ email/SMS/push/etc  │
                         └─────────────────────┘
```

---

# 5. **Core Components**

---

## 5.1 **API Gateway**

Gateway is the public-facing entrypoint.
Responsibilities:

* REST → gRPC translation
* Authentication proxy
* Session verification
* Rate limiting
* Request signing
* Security mode enforcement
* Input validation
* Multi-tenant routing

Supported security:

* **HMAC** — time-limited request signatures
* **Static Sign**
* **mTLS** — client certificate authentication

Gateway ensures **all traffic inside the platform is trusted and uniform**.

---

## 5.2 **MasterService**

MasterService is the **Control Plane** of Gufo.

Responsibilities:

### ✔ Microservice Registry

Stores metadata of all microservices:

* name
* host/port
* groups
* internal/external flags
* version
* master/replica info

### ✔ Routing Logic

Decides where to send internal requests (service discovery).

### ✔ Cron/Job Coordination

Guarantees that **only one replica** of a service executes scheduled jobs.

### ✔ EntryPoint and Version Sync

Ensures microservices load correct logic/structural versions.

### ✔ Health Checks

Local `/health` + registry health states.

MasterService **must be isolated**, like DNS or Kubernetes API server.

---

## 5.3 **Session Service**

SessionService stores user sessions in local **Redis**.

### Why separate?

* Sessions are checked **on every request**
* Sessions must be **geo-local**
* Each region has its own "session cluster"
* Millions of concurrent users → millions of hashes → global session DB impossible
* Scaling by region is trivial:

    * redis-eu
    * redis-africa
    * redis-asia

Auth should *never* handle that load → correct separation.

---

## 5.4 **Auth Service**

Auth = **CPU heavy** service:

* password hashing
* verification
* MFA
* brute-force defense
* login limits
* token generation
* audit events

Auth is a **global service**, unlike sessions.
It must NOT be merged with sessions — different scaling profiles.

---

## 5.5 **Registration Service (Reg)**

Registration workflows differ dramatically between clients:

* simple email + password
* multi-step forms
* KYC
* SMS verification
* invitation-based registration
* admin-created accounts
* webhook-based registration
* custom provider integrations

Reg must be **customizable**, sometimes replaced entirely.

Therefore it **must be a separate microservice**.

---

## 5.6 **Notifications Service**

Notifications vary by:

* geography
* industry
* message type
* provider
* reliability
* volume

Supported channels:

* email (SMTP/SES/etc.)
* SMS
* Telegram
* WhatsApp
* Push
* WebPush
* Slack
* Custom providers

Notifications = transport layer → should not live inside Auth or any other service.

---

## 5.7 **Rights Service**

Rights & permissions is a **domain in itself**:

* RBAC
* ABAC
* ACL
* Tenant isolation
* Feature flagging
* Resource-level permissions
* Context-based permissions

Rights must be independent, since permission logic is:

* complex
* evolving
* business-specific
* often integrated with external systems

---

# 6. **Internal Communication Model**

### Protocol:

**Gufo Internal gRPC Protocol (Reverse API)**
Request object: `pb.Request`
Response object: `pb.Response`

### Guarantees:

* authenticated
* signed
* validated
* uniform for all microservices
* versioned
* backward compatible

This provides a **stable internal API** independent of REST schemas.

---

# 7. **Security Architecture**

Gufo uses **three-layer security**:

### Layer 1 — Transport

* TLS
* mTLS (optional)

### Layer 2 — Message Integrity

* HMAC signatures
* static Sign tokens
* expiration timestamps
* replay protection

### Layer 3 — Application Layer

* sessions
* rights
* tokens
* rate limits

MasterService enforces `security.mode` globally.

---

# 8. **Data Architecture**

Each service owns its own storage:

* Session → Redis
* Auth → SQL/NoSQL
* Rights → SQL/graph DB
* Notifications → internal queue + SMTP/SMS providers
* Registration → SQL (if needed)

There is **no global database**.
Gufo follows the principle: **shared nothing**.

---

# 9. **Geo-Distributed Model**

A core part of Gufo design.

### Regions contain:

* local SessionService
* local Notification providers
* optionally local Auth replicas for read-heavy operations

But have:

* global identity (User)
* global Rights
* global MasterService network (clustered)

Example:

```
(EU)    ── master-eu     ── session-eu     ── notifications-eu
(US)    ── master-us     ── session-us     ── notifications-us
(AFR)   ── master-africa ── session-africa ── notifications-africa
(ASIA)  ── master-asia   ── session-asia   ── notifications-asia
```

---

# 10. **Scaling Strategy**

### Horizontally scalable components:

* Session (Redis replication)
* Notifications
* Auth
* Rights (depends on backend)
* MasterService (active-active with gossip or LB mode)

### Vertically scalable components:

* Auth (CPU-bound hashing)
* Rights (complex queries)

---

# 11. **Deployment Modes**

Gufo supports:

### ✔ Local (single-host)

For development/testing.

### ✔ Standard

MasterService + Gateway + grouped microservices.

### ✔ Enterprise

All 6 microservices separate.

### ✔ Geo-distributed

Several masterservice clusters + region-local sessions.

### ✔ Embedded

Only Gateway + MasterService + basic Auth (for small projects).

---

# 12. **Extensibility Model**

Gufo is built as a **plug-and-play platform**:

You can replace:

* Auth → your own auth service
* Notifications → your own SMS/email provider
* Registration → custom flow
* Rights → external RBAC engine
* Session → any Redis cluster

You can add:

* billing
* CRM
* analytics
* catalog
* storage
* payment services

Everything talks through the same unified protocol.

---

# 13. **Glossary**

* **Control Plane** — MasterService & configuration.
* **Data Plane** — Session, Auth, Rights.
* **Gufo Internal Protocol** — unified gRPC contract.
* **Reverse API** — bidirectional streaming in Gufo.
* **Registry** — list of microservices + routing metadata.
* **Session** — short-lived state of a logged user.
* **Rights** — access control logic.
* **Notifications** — messaging transport layer.




