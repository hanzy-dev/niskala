# Engineering Decisions

## 1. Backend-heavy by design
Niskala is intentionally backend-heavy because the main value of the project is in transactional correctness, reliability, and system behavior rather than frontend complexity.

## 2. Go for the main API
Go is used as the orchestration layer because it is well-suited for REST APIs, transaction-heavy backend workflows, and operationally simple services.

## 3. Gin instead of a lighter router
Gin is chosen to better align with common backend job requirements and provide a familiar, productive HTTP framework for the main API.

## 4. Rust as an isolated pricing service
Rust is used for a small, isolated pricing service to demonstrate clear service boundaries and dependency failure handling. It does not own the main database.

## 5. REST between Go and Rust
Internal service-to-service communication uses HTTP/JSON instead of gRPC to keep the project easier to build, test, and debug while preserving the core engineering story.

## 6. Postgres as source of truth
All transactional state is stored in Postgres. Product prices are re-read during checkout and order item snapshots are persisted for historical correctness.

## 7. Redis only for product reads
Redis is used as a cache layer for read-heavy catalog access. It is not used as the system of record.

## 8. Idempotency is mandatory
Checkout uses idempotency keys to prevent duplicate order creation during retries or repeated submissions.

## 9. Stock consistency matters more than raw throughput
The checkout flow uses transaction boundaries and ordered row locking to keep stock updates consistent under concurrent requests.

## 10. Pricing availability must not block checkout
If the pricing service is unavailable, checkout continues with normal pricing and explicitly reports fallback usage.

## 11. Email is post-commit and best-effort
Order confirmation email is sent after commit. Email delivery failure must not roll back the order.

## 12. Observability is part of the design
Structured logs, consistent error formats, and correlation IDs are included because debugging system behavior is part of backend engineering maturity.