# Niskala

A backend-heavy commerce system focused on checkout correctness, idempotent order processing, stock-safe transactions, and resilience against downstream service failures.

## Overview

Niskala is a backend-heavy commerce project designed to demonstrate:
- safe checkout flow design
- idempotent order creation
- stock-safe transaction handling
- pricing service fallback
- Redis-backed local infrastructure
- cross-service observability

It is intentionally built to emphasize backend engineering depth over UI complexity.

## Stack

### Frontend
- Vue 3
- TypeScript
- Pinia
- Vue Router
- Axios

### API
- Go
- Gin

### Pricing Service
- Rust
- Axum

### Data and Infra
- Supabase Auth
- Supabase Postgres
- Redis
- Docker Compose
- GitHub Actions

## Current capabilities

- product list and product detail flow
- cart add/update/delete flow
- checkout placeholder flow
- order history and order detail flow
- idempotency key replay support
- pricing service integration
- pricing fallback foundation
- correlation ID propagation
- consistent error response format
- local infra with Redis

## Project structure

- `apps/web` — Vue frontend
- `apps/api` — Go API
- `apps/pricing` — Rust pricing service
- `docs` — architecture and scope docs
- `scripts` — local helper scripts

## Run locally

### Web
```bash
cd apps/web
npm install
npm run dev
```

### API

```
cd apps/api
go run ./cmd/server
```

### Pricing service

```
cd apps/pricing
cargo run
```

### Redis

```
docker compose up -d
```

## Tests

### API

```
cd apps/api
go test ./...
```

### Pricing

```
cd apps/pricing
cargo test
```

## Engineering notes

- money values use integer cents
- checkout requires Idempotency-Key
- pricing service failures can fall back to normal pricing
- correlation IDs are propagated through HTTP requests
- current auth flow uses debug headers as a local development placeholder

## Next evolution

- replace debug auth with Supabase JWT verification
- move in-memory flows to Postgres-backed repositories
- add transactional checkout with row locking
- persist idempotency records in the database
- add Redis-backed product caching
- enrich structured logging and metrics