# Niskala

A high-integrity commerce engine designed for critical transaction safety. Niskala prioritizes checkout correctness and system resilience, ensuring zero-stock leaks and idempotent processing even when downstream services fail.

## Stack

- Frontend: Vue 3 + TypeScript + Pinia + Vite
- API: Go + Gin
- Pricing service: Rust
- Database: Postgres (Supabase)
- Cache: Redis
- Auth: Supabase Auth (Google OAuth)

## Current Capabilities

### User flow
- Login with Google via Supabase
- Browse products from real backend
- View product detail
- Add items to cart
- Checkout with idempotency key
- View orders and order details

### Admin flow
- Protected admin routes
- Product list for admin
- Stock update from frontend
- Create product from frontend

### Backend guarantees
- Products, carts, orders, order items persisted in Postgres
- Membership-based admin authorization
- Redis caching for product reads
- Pricing service integration with fallback mode
- Transactional checkout with stock locking
- Idempotent checkout persisted in database

## Project Structure

- `apps/web` — frontend app
- `apps/api` — Go API
- `apps/pricing` — Rust pricing service
- `docs` — architecture and scope docs
- `scripts` — local development helpers

## Local Development

### 1. Start Redis
```bash
docker compose up -d
```

2. Start API
```bash
cd apps/api
go run ./cmd/server
```

3. Start pricing service
```bash
cd apps/pricing
cargo run
```

4. Start frontend
```bash
cd apps/web
npm run dev
```

## Environment Files
### API

Copy:

- `apps/api/.env.example` → `apps/api/.env`

### Web

Copy:

- `apps/web/.env.example` → `apps/web/.env`

### Notes

- Supabase project is used for Postgres and Auth
- Google OAuth is configured through Supabase Auth
- Admin access depends on the memberships table role