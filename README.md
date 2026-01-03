## Tech Stack

- Golang 1.25.3
- Chi - routing
- SQLC - queries
- Goose - migrations
- Docker

## Architecture Overview

- Transport layer - http/grpc - cmd/api.go, cmd/main.go
- Service layer - all business logic ex: getOrders, getOrderItems
- Repository layer - postgres - data fetching - SELECT \* FROM products;
- CMD folder - used for executables, CLI, etc.
- Internal folder - reserved folder that doesn't export.

## Helpful Dev Commands

- `docker compose up -d` - start the postgres database
- `go run cmd/*.go` - run the application
- `sqlc generate` - will regenerate your SQL queries
- `goose -s create create_products sql` - will create migrations under `internal/adapter/postgres/migrations`
- `goose up` - apply all available migrations
- `goose down` - roll back a single migration from the current version

## Next Steps

### 1. **Order Item Status Enhancements** (DONE)

**Objective:**  
Support partial fulfillment by introducing per-item order statuses.

- Implement an `order_item.status` field to represent the fulfillment state of each item within an order.
- Enable the following logic:
  - For multi-item orders, if any item is out of stock, set its status to **`backordered`**.
  - Items that are in stock should have their status set to **`fulfilled`**.
- **Example:**
  - An order contains Item A (in stock) and Item B (out of stock):
    - Item A → `fulfilled`
    - Item B → `backordered`

### 2. **Preventing Overselling & Handling Stock Race Conditions**

**Current Challenge:**  
Inventory levels are updated per item immediately, which allows concurrent orders for the same product to both pass stock checks. This can result in overselling.

**Recommended Solutions:**

#### a) Pessimistic Locking

- Lock inventory rows using `SELECT FOR UPDATE` within a transaction.
- Implement a timeout (e.g., 2 seconds) to minimize the risk of long-held locks.
- _**Pros:** Provides strong consistency and is effective for moderate traffic volumes._
- _**Cons:** May cause contention for popular items under high load._

#### b) Redis-based Inventory Reservation

- Upon ordering, create a temporary reservation of inventory in Redis.
- Set a reservation expiry (e.g., 15 minutes) to auto-release stock in case of order abandonment.
- Optionally, use Redis counters for tracking "hot items" and supporting high-demand sales scenarios.
- Supplement with background jobs for clean-up or reconciliation if necessary.
- _**Pros:** Fast, scalable, and suitable for events like flash sales._
- _**Cons:** Additional complexity; requires operational Redis._

#### c) Job Queue for Order Processing

- Queue incoming order requests to serialize access to inventory for popular or high-traffic items.
- Ensures consistent stock updates and avoids race conditions under high concurrency.
- _**Pros:** Handles conflict resolution at scale and provides graceful degradation._
- _**Cons:** May introduce order processing latency._

---
