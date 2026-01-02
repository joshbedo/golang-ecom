Tech Stack

- Golang 1.25.3
- Chi - routing
- SqlC - queries
- Goose - migrations

Architecture Overview

- Transport layer - http/grpc - cmd/api.go, cmd/main.go
- Service layer - all business logic ex: getOrders, getOrderItems
- Repository layer - postgres - data fetching - SELECT \* FROM products;
- CMD folder - used for executables, CLI, etc.
- Internal folder - reserved folder that doesn't export.

Helpful Dev Commands

- `go run cmd/*.go` - run the application
- `sqlc generate` - will regenerate your SQL queries
- `goose -s create create_products sql` - will create your migrations
- `goose up` - apply all available migrations
- `goose down` - roll back a single migration from the current version
