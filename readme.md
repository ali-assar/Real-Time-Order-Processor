## Real-Time Order Processor

This is a backend service that simulates processing thousands of incoming e-commerce orders per second.

### Project Setup

order-processor/
 ├─ cmd/
 │   └─ server/        # main.go entrypoint
 ├─ internal/
 │   ├─ handler/       # HTTP handlers
 │   ├─ processor/     # business logic
 │   └─ storage/       # in-memory or database
 ├─ go.mod
 └─ README.md

