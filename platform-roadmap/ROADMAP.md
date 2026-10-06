# Infrastructure Software Roadmap

**Goal:** Web developer → infrastructure / platform software engineer: a
high-paying role that is mostly writing systems code (Kubernetes controllers,
developer tooling, networking and storage software), with Go as the main
language and a low-level stretch (Rust or C++) later.

**Background:** ~1 year of C at École 42 (pointers, memory, Unix basics), plus
web development. Go work started 2026-08-19.

**Start date:** 2026-08-19 · **Target:** no fixed date. Timeline follows real
weekly hours (see "Pacing").

**How progress works:** tick a box only when the thing exists and works —
code runs, tests pass, or you can explain it without notes. Reading about
something doesn't count. If you've *used* something without understanding it,
leave it unticked with a short note.

---

## Pacing

Original plan assumed ~34 h/week. Actual hours have been well under that and
are rising. Rough effort per phase is below (estimates, to be re-baselined
after ~2 weeks of logged hours); calendar time = hours ÷ your weekly hours.

| Phase | Est. hours |
|---|---|
| 1 Go fluency (remaining) | ~20 |
| 2 Backend essentials | ~60 |
| 3 Linux and networking | ~120 |
| 4 Containers and Kubernetes | ~150 |
| 5 Infrastructure software builds | ~150 |
| 6 IaC, CI/CD, cloud | ~100 |
| 7 Observability and reliability | ~80 |
| 8 Distributed systems | ~80 |
| 9 Flagship | ~200 |

Total ≈ 960 h. Recovery: take a lighter week every 6–8 weeks.

---

## Phase 1 — Go fluency (mostly done)

Concepts
- [x] Syntax, types, structs, methods, slices, maps
- [x] Interfaces and embedding — _interfaces explained 2026-09-30; embedding not covered yet_
- [x] Error handling (wrapping, `errors.Is/As`, sentinel errors) — _explained 2026-09-30; custom error types not covered yet_
- [x] Testing (table tests, `go test`, `-race`, `-cover`) — _explained 2026-10-06; `t.Parallel`, `httptest`, benchmarks not covered yet_
- [x] Modules and package layout — _explained 2026-10-06; `go get`, `internal/`, workspaces not covered yet_
- [x] Goroutines and channels — _explained 2026-09-30; `select` patterns, fan-in/out not covered yet_
- [x] Mutexes and `sync` package — _explained 2026-09-30; `RWMutex`, `sync.Once`, `atomic` not covered yet_
- [x] `context` (cancellation, timeouts) — _cancellation explained 2026-09-30; `WithTimeout` and `ctx` in HTTP calls not covered yet_

Builds
- [x] A CLI tool (flags, stdin/stdout, exit codes) — `week01/todo-cli`
- [x] Concurrent downloader (worker goroutines, cancellation, errors) — `week02/concurrent-url-checker`
- [x] Small HTTP API with tests — `week02/todo-http-api`

**Understanding pass:** before Phase 2, explain each unticked concept above
without notes. The agent writes the `concepts/` entry from your explanation.

## Phase 2 — Backend essentials (short; a means, not the goal)

Concepts
- [ ] HTTP in depth (methods, status codes, headers, `net/http`)
- [ ] Middleware (logging, recovery, auth)
- [ ] Authentication (sessions vs tokens, JWT, password hashing)
- [ ] PostgreSQL basics and SQL
- [ ] Indexing
- [ ] Transactions and isolation levels

Builds
- [ ] API with Postgres persistence, auth, middleware and tests

## Phase 3 — Linux and networking (your strength; build on it)

Concepts
- [ ] Processes, threads, scheduling — _in flight: `week03/process-lifecycle`_
- [ ] Signals, permissions, systemd, logs
- [ ] Virtual memory
- [ ] File descriptors, epoll / I/O models
- [ ] cgroups and namespaces (what a container actually is)
- [ ] TCP vs UDP, sockets — _in flight: `week03/tcp-echo-server`_
- [ ] HTTP, DNS, TLS
- [ ] Load balancing and proxying
- [ ] Concurrency patterns (pipelines, fan-in/out, worker pools)
- [ ] Profiling (`pprof`) and benchmarking
- [ ] Debugging tools: `strace`, `ss`, `tcpdump`

Builds
- [ ] Worker pool with graceful shutdown
- [ ] TCP server beyond echo (framing, concurrency, timeouts)
- [ ] Rate limiter (token bucket or similar), benchmarked
- [ ] A tiny container runtime using namespaces and cgroups (Go)

## Phase 4 — Containers and Kubernetes as a developer

Concepts
- [ ] Docker (images, containers, volumes, compose, multi-stage builds)
- [ ] Kubernetes core objects (pods, deployments, services, config, RBAC)
- [ ] Running a real cluster locally (kind or k3s)
- [ ] `client-go` and the Kubernetes API
- [ ] The reconciliation loop and controller pattern

Builds
- [ ] Deploy your Phase 2 API on a local cluster
- [ ] A Kubernetes controller or operator in Go (custom resource plus reconciler)

## Phase 5 — Infrastructure software builds

Build at least two (tick what you build), designed as tools other engineers use:
- [ ] Deployment CLI
- [ ] Job scheduler or queue
- [ ] Configuration or feature-flag service
- [ ] API gateway or reverse proxy
- [ ] Service-discovery prototype
- [ ] Key-value store with a write-ahead log

**Checkpoint (end of Phase 5):** revisit language choice. If the low-level side
is the draw, pick Rust or C++ for the next build. Also decide the target role
family. Start reading job postings seriously.

## Phase 6 — Infrastructure as code, CI/CD, cloud (light but real)

- [ ] Terraform: provision something real
- [ ] CI/CD pipeline for one of your projects
- [ ] One cloud (AWS or GCP): IAM, networking, managed Kubernetes basics
- [ ] Deploy one project to the cloud with IaC
- [ ] Secrets and least-privilege basics

## Phase 7 — Observability and reliability

- [ ] Structured logging
- [ ] Metrics (Prometheus-style), exposed from your own service
- [ ] Tracing
- [ ] SLOs and alerting for one service
- [ ] Break something on purpose and write a short incident write-up

## Phase 8 — Distributed systems

Concepts
- [ ] Consistency models
- [ ] Replication
- [ ] Caching strategies and invalidation
- [ ] Queues and delivery guarantees
- [ ] Retries, backoff, idempotency
- [ ] Failure modes and tradeoffs (CAP, partitions, timeouts)

Builds
- [ ] Toy replicated key-value store, or caching + retries + idempotency added to an earlier project

## Phase 9 — Flagship project

Internal developer platform: a developer runs one command or API call and gets
a service deployed, observable and running on Kubernetes.
- [ ] Authentication
- [ ] Worker queues
- [ ] Metrics
- [ ] Logging
- [ ] PostgreSQL
- [ ] Docker and Kubernetes
- [ ] Deployed with IaC and CI/CD
- [ ] Tests
- [ ] Documentation (README, architecture notes)

## Ongoing

Read excellent code — log each project you study (package layout, interfaces,
tests, errors). Good candidates: Kubernetes controllers, etcd, Terraform providers.
- [ ] Project 1: _______
- [ ] Project 2: _______
- [ ] Project 3: _______

Low-level stretch (Rust or C++), after Phase 5:
- [ ] One small systems project in Rust or C++

Job search:
- [ ] Mid-plan: apply to adjacent roles (backend, infra tooling, DevOps/SRE) as calibration, around the end of Phase 5–6
- [ ] Portfolio README for each flagship-level project

---

## Bonus hours priority

1. Improve existing projects
2. Read source code
3. Refactor and polish
4. Write technical notes (→ `concepts/`)
5. Linux mastery
6. Profiling and performance work

## Certifications

Portfolio first. Certifications can help pass recruiter filters: consider AWS
Solutions Architect or CKA around Phases 4–6, only after hands-on experience.

---

## Adopted adjustments

Adopted 2026-09-30, after a review of the original AI-generated plan and
current job-market data:
- Re-baselined capacity to real hours (see "Pacing"); dropped the fixed 10–12 month target.
- Understanding pass on Phase 1 before continuing.
- Retargeted from generic backend/platform to infrastructure software in Go.
- Linux and networking deepened; Kubernetes moved earlier and made developer-focused (`client-go`, controllers).
- IaC, CI/CD and cloud given a real (but light) phase.
- Observability and reliability practice added; a Phase 8 build added.
- Mid-plan job-search checkpoint and Rust/C++ stretch added.

## Proposed adjustments (not adopted yet — decide when they come up)

- **Recovery weeks.** A lighter week every 6–8 weeks. Raise around week 6 of steady hours.
- **Phase 3 buffer.** Most ambitious phase; allow it to run long rather than rushing it. Raise if it starts slipping.
- **Domain choice.** At the Phase 5 checkpoint, decide between infrastructure/platform, systems/OS-adjacent, performance-critical or embedded, and pick Rust or C++ accordingly.
