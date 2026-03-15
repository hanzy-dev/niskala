# Architecture

## Overview

Niskala is a backend-heavy commerce system with three main application components:

1. **Web app** built with Vue 3
2. **Main API** built with Go and Gin
3. **Pricing service** built with Rust and Axum

Supporting infrastructure:
- Supabase Auth for authentication
- Supabase Postgres as the primary database
- Redis for product read caching
- Resend for best-effort order confirmation emails

## High-level component roles

### Web app
The frontend is a thin client layer responsible for:
- login flow
- product browsing
- cart interaction
- checkout initiation
- displaying order history
- showing degraded mode/fallback messages

The frontend is **not trusted** for:
- final prices
- payment totals
- admin authorization
- discount decisions

### Main API
The Go API is the transactional core of the system.

It is responsible for:
- verifying Supabase JWTs
- loading user identity
- enforcing membership-based RBAC
- serving product and cart APIs
- orchestrating checkout
- enforcing idempotency
- locking stock rows safely
- integrating with the pricing service
- applying fallback when pricing is unavailable
- writing orders and order item snapshots
- emitting logs and health signals

### Pricing service
The Rust service is isolated from the main database.

It is responsible for:
- calculating subtotal
- calculating discount
- calculating total
- returning a pricing breakdown

The pricing service only receives item snapshots from Go and returns a deterministic pricing result.

## Communication model

### Frontend to API
- REST over HTTP
- JWT-based authenticated requests
- correlation ID propagation
- idempotency key on checkout

### API to pricing service
- internal REST over HTTP/JSON
- timeout and retry strategy
- fallback if service is unavailable

## Data ownership

### Postgres is the source of truth for:
- products
- carts
- cart_items
- orders
- order_items
- memberships
- idempotency_keys

### Redis is only used for:
- product list cache
- product detail cache

Redis is not the source of truth.

## Key system properties

- checkout is transaction-based
- order creation is idempotent
- stock updates are concurrency-aware
- pricing dependency failure should not block checkout
- all money values use integer cents
- request flow is traceable through correlation IDs