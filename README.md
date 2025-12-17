# Distributed Rate Limiter (Go + Redis)

A production-grade **distributed rate limiter** built in Go using Redis and Lua, supporting:
- Token Bucket algorithm
- Per-IP and per-user limits
- Fail-open and fail-closed strategies
- Local in-memory fallback during Redis outages
- Prometheus metrics for observability

---

## 🚀 Features

- **Distributed rate limiting** using Redis + Lua (atomic operations)
- **Token Bucket algorithm**
- **Per-identity limits**
  - Anonymous users → IP-based
  - Authenticated users → User-based
- **Failure strategies**
  - Fail-Open (allow traffic with fallback limiter)
  - Fail-Closed (block traffic)
- **Local fallback limiter** when Redis is unavailable
- **Prometheus metrics** exposed via `/metrics`
- **IPv4/IPv6 normalization**
- Clean middleware-based architecture

---

## 🧠 High-Level Architecture

Client
|
v
HTTP Server (net/http)
|
v
Rate Limit Middleware
|
+-- Redis Available?
| |
| +-- YES --> Redis + Lua (Token Bucket)
| |
| +-- NO
| |
| +-- Fail-Open --> Local In-Memory Limiter
| |
| +-- Fail-Closed --> Reject Request (429)
|
v
Application Handler



---

## 🏗️ Tech Stack

- **Go** (net/http)
- **Redis** (distributed state)
- **Lua** (atomic token bucket logic)
- **Prometheus** (metrics & observability)
- **Docker** (Redis container)

---

## 🔐 Rate Limiting Logic

### Token Bucket Parameters
- `capacity` → max burst size
- `refillRate` → tokens per second

### Identity Strategy
| Request Type | Key Used |
|-------------|----------|
Anonymous | `ip:<client_ip>` |
Authenticated | `user:<user_id>` |

---

## ⚠️ Failure Strategies

### Fail-Open
- Redis down → allow traffic
- Use local in-memory limiter
- Adds header: `X-RateLimit-Degraded: true`

### Fail-Closed
- Redis down → block all traffic
- Returns `429 Too Many Requests`

Configured via `config.FailureStrategy`.

---

## 📊 Metrics (Prometheus)

Exposed at:
GET /metrics



### Available Metrics
- `rate_limit_allowed_total{identity="ip|user"}`
- `rate_limit_blocked_total{identity="ip|user"}`
- `rate_limit_degraded_total`
- `rate_limit_errors_total`

---

## 🧪 Running Locally

### 1️⃣ Start Redis
```bash
docker run -d -p 6379:6379 redis
2️⃣ Run Server

go run cmd/server/main.go
3️⃣ Test

curl http://localhost:8080/health
curl http://localhost:8080/metrics
📁 Project Structure

cmd/server        → entry point
internal/config   → configuration
internal/limiter  → Redis + Lua rate limiter
internal/middleware → HTTP middleware
internal/metrics  → Prometheus metrics
pkg/response      → shared response helpers
🎯 Why This Project?
This project demonstrates:

Distributed systems thinking

Failure handling and graceful degradation

Middleware design in Go

Observability and production readiness

📌 Future Improvements
JWT validation

Sliding window algorithm

Dynamic per-tier limits

Grafana dashboards

Kubernetes deployment

👤 Author
Paras Tayade
Backend / Systems Engineer



---


