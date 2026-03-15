# Niskala

A backend-heavy commerce system focused on checkout correctness, idempotent order processing, stock-safe transactions, and resilience against downstream service failures.

## Why this project exists

Niskala is built to demonstrate backend engineering maturity through a commerce workflow that prioritizes transaction correctness, retry safety, stock consistency, resilience, and observability.

Instead of focusing on feature breadth, this project focuses on engineering depth in the parts of commerce systems that commonly fail in real-world environments:
- duplicate checkout
- stock inconsistency
- race conditions
- trust boundary violations
- downstream dependency failure
- poor request traceability

## Core goals

- Design a safe transactional checkout flow
- Prevent duplicate orders with idempotency keys
- Keep stock consistent during concurrent checkout
- Continue checkout when pricing service fails
- Trace requests across frontend, API, and pricing service
- Keep architecture clean, small, and reviewable

## Final stack

### Frontend
- Vue 3
- TypeScript
- Pinia
- Vue Router

### Main API
- Go
- Gin

### Pricing Service
- Rust
- Axum

### Data and Infra
- Supabase Auth
- Supabase Postgres
- Redis
- Resend
- Docker Compose
- GitHub Actions

## Project shape

Niskala is intentionally built as a backend-heavy system.

The frontend exists to demonstrate the product flow:
- login
- product catalog
- product detail
- cart
- checkout
- order history
- admin product management

The main value of the project is in the backend:
- transactional checkout
- idempotency
- stock locking
- pricing fallback
- cache strategy
- observability
- clean service boundaries

## Documentation

- [Architecture](./docs/architecture.md)
- [Scope](./docs/scope.md)
- [Engineering decisions](./docs/engineering-decisions.md)

## Status

This project is being built incrementally in small, reviewable batches with clean commit history.