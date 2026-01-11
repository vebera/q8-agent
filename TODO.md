# Q8 Agent Implementation Plan

This agent manages tenant environments on host servers, providing a secure API for the Q8 Main Server to coordinate infrastructure without SSH.

## Architecture
- **Language**: Golang
- **Style**: Uber Go Style Guide
- **Communication**: HTTP REST API (private network)
- **Authentication**: Bearer Token
- **Docker**: CLI-wrapped interaction (mapping `/var/run/docker.sock`)

## Directory Structure
```text
q8-agent/
├── cmd/agent/          # Entry point
├── internal/
│   ├── api/            # HTTP handlers and middleware
│   ├── domain/         # Domain types (Tenant, Request/Response models)
│   ├── config/         # Environment-based configuration
│   ├── docker/         # Docker Compose execution engine
│   ├── fs/             # Tenant directory and config file manager
│   └── service/        # Business logic & orchestration
├── Dockerfile          # Multi-stage build (Alpine + Docker CLI)
└── TODO.md             # Implementation roadmap
```

## Phase 1: Foundation ✅
- [x] Initialize Go module and package structure.
- [x] Implement `internal/config` with env-var support (`Q8_AGENT_PORT`, `Q8_AGENT_ADMIN_TOKEN`, etc.).
- [x] Implement `internal/api/middleware` for Bearer Token authentication.
- [x] Setup basic `http.ServeMux` with structured health checks.

## Phase 2: Core Orchestration ✅
- [x] Implement `FSManager` for safe tenant workspace management under `/opt/tenants`.
- [x] Implement `DockerRunner` for atomic `docker compose` operations.
- [x] Support project isolation using namespaces (`q8-{subdomain}`).
- [x] Implement `Orchestrator` service to chain FS and Docker operations.

## Phase 3: APIs & Telemetry ⏳
- [x] `POST /v1/tenants/provision`: full setup (mkdir + write configs + pull + up).
- [x] `POST /v1/tenants/teardown/{subdomain}`: graceful shutdown + file cleanup.
- [x] `POST /v1/tenants/restart/{subdomain}`: restart containers.
- [x] `GET /v1/tenants/status/{subdomain}`: JSON status of all services in stack.
- [x] **NEW**: `GET /v1/tenants/logs/{subdomain}`: Stream or tail logs from specific/all services.
- [x] **NEW**: `GET /v1/tenants/images/{subdomain}`: Report current image IDs and tags running.
- [ ] **NEW**: `GET /v1/system/stats`: Host-level telemetry (CPU/RAM/Disk) for load balancing by Main Server.
- [ ] **NEW**: `POST /v1/tenants/update`: Lightweight image update (pull + up) without full re-provisioning.

## Phase 4: Production Readiness & Security 🛠️
- [x] Multi-stage `Dockerfile` with Docker-CLI-Compose support.
- [ ] Implement advanced cleanup logic (pruning orphan volumes/networks per tenant).
- [ ] Add support for `DOCKER_REGISTRY` credentials (auth against private registries).
- [ ] **Pre-flight checks**: Implement port availability validation before starting containers.
- [ ] **Resource Limits**: Configurable max tenants per agent.
- [ ] **Concurrent Safety**: Mutex-protected operations per tenant to prevent race conditions during updates.

## Phase 5: Testing & Integration 🧪
- [ ] **Unit Testing**: Implement table-driven tests for FS and Config packages.
- [ ] **Mocking**: Add interface-based mocks for Docker and FS to test Orchestrator logic.
- [ ] **Integration Test**: Create a script simulating the Main Server lifecycle (Provision -> Status -> Update -> Teardown).
- [ ] **Main Server Integration**: Refine `SSHService` in `v0-qate-landing` to support a "Local Agent" mode.
