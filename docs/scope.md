# Final Scope

## Project positioning

Niskala is a backend-heavy commerce project built to demonstrate:
- checkout correctness
- idempotent order processing
- concurrency-safe stock handling
- resilient service integration
- Redis-backed caching
- cross-service observability

## User-facing scope

- Google login via Supabase
- product catalog
- product detail
- simple search
- cart
- checkout
- checkout success
- order history
- order detail

## Admin scope

- create product
- update product
- update stock
- membership-based admin authorization

## System scope

- transactional checkout
- idempotency key handling
- stock row locking
- consistent lock ordering
- Rust pricing service integration
- pricing fallback
- Redis cache for product reads
- structured logs
- correlation ID propagation
- consistent error response format
- health endpoint
- async email after commit
- tests
- Docker Compose
- GitHub Actions

## Explicitly out of scope

These are intentionally excluded to keep the project focused:
- gRPC
- Kafka or RabbitMQ
- complex payment gateway integration
- shipment tracking
- multi-vendor marketplace flows
- reviews and ratings
- loyalty systems
- large analytics features
- extra microservices beyond the pricing service