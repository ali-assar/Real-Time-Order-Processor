# Real-Time Order Processor

A high-performance, concurrent backend service designed to process thousands of e-commerce orders per second with real-time processing capabilities, priority-based queuing, and comprehensive monitoring.

## 🚀 Features

- **Concurrent Processing**: Worker pool architecture for high-throughput order processing
- **Priority-Based Queuing**: Orders processed based on priority levels (High/Medium/Low)
- **Real-Time Monitoring**: Live statistics and health monitoring endpoints
- **Business Logic Processing**: Intelligent order validation and business rule application
- **Error Handling**: Comprehensive error handling with detailed logging
- **Scalable Architecture**: Configurable worker pool size and buffer capacity

## 🏗️ Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   HTTP API  │───▶│ Order Pool │───▶│  Workers   │
│   (Port 8080)│    │   (Buffer)  │    │ (Configurable)│
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                           ▼                   ▼
                    ┌─────────────┐    ┌─────────────┐
                    │   Results   │    │   Logging   │
                    │   Channel   │    │  & Metrics  │
                    └─────────────┘    └─────────────┘
```

## 📁 Project Structure

```
Real-Time-Order-Processor/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── handler/             # HTTP request handlers
│   │   ├── handler.go       # Order creation and processing handlers
│   │   └── router.go        # Route registration
│   ├── processor/           # Business logic and worker pool
│   │   └── pool.go          # Worker pool implementation
│   └── pkg/
│       └── models/          # Data models and validation
│           └── model.go     # Order, ProcessedOrder, and stats models
├── go.mod                   # Go module dependencies
└── README.md               # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24.4 or higher
- Git

### Installation & Running

1. **Clone the repository**
   ```bash
   git clone https://github.com/ali-assar/Real-Time-Order-Processor.git
   cd Real-Time-Order-Processor
   ```

2. **Run the service**
   ```bash
   go run cmd/main.go
   ```

3. **Service will start on port 8080**
   ```
   API listening on :8080
   ```

## 📡 API Endpoints

### 1. Create Order
**POST** `/orders`

Creates a new order and queues it for processing.

**Request Body:**
```json
{
  "id": "order_123",
  "amount": 99.99,
  "items": ["item1", "item2"],
  "customer": "john.doe@example.com",
  "status": "pending",
  "address": "123 Main St, City, Country",
  "priority": 1,
  "notes": "Handle with care"
}
```

**Priority Levels:**
- `1` = High Priority (processed first)
- `2` = Medium Priority (default)
- `3` = Low Priority

**Response:**
```json
{
  "id": "order_123",
  "amount": 99.99,
  "items": ["item1", "item2"],
  "customer": "john.doe@example.com",
  "status": "pending",
  "created_at": "2024-01-15T10:30:00Z",
  "address": "123 Main St, City, Country",
  "priority": 1,
  "notes": "Handle with care"
}
```

### 2. Get Processing Statistics
**GET** `/stats`

Returns real-time processing statistics.

**Response:**
```json
{
  "total_processed": 150,
  "success_count": 145,
  "error_count": 5,
  "average_process_time_ms": 45.2,
  "active_workers": 10,
  "queue_length": 3,
  "uptime_seconds": 3600
}
```

### 3. Health Check
**GET** `/health`

Returns the health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1705312200,
  "pool": {
    "healthy": true,
    "queue_length": 3,
    "workers": 10
  }
}
```

## ⚙️ Configuration

The service can be configured by modifying the following parameters in `cmd/main.go`:

```go
pool := processor.Start(context.Background(), 10, 100)
//                              │        │
//                              │        └── Buffer size (orders)
//                              └── Number of workers
```

**Recommended configurations:**
- **Development**: 5 workers, 50 buffer
- **Production**: 20-50 workers, 1000+ buffer
- **High-load**: 100+ workers, 5000+ buffer

## 🔧 Business Logic

### Order Processing Flow

1. **Validation**: Order data validation and business rule checks
2. **Priority Processing**: Orders processed based on priority level
3. **Business Rules**: 
   - Orders > $1000 marked for priority processing
   - High priority orders expedited
   - Amount limits enforced ($10,000 max)
   - Item count limits (50 items max)
4. **Result Generation**: Processing results with timing and worker information

### Processing States

- `pending` → Initial state
- `processing` → Being processed
- `priority_processing` → High-value order
- `expedited` → High priority order
- `completed` → Successfully processed

## 📊 Monitoring & Metrics

The service provides comprehensive monitoring:

- **Real-time Statistics**: Processing counts, success rates, timing
- **Health Monitoring**: Service health and pool status
- **Performance Metrics**: Average processing time, queue length
- **Worker Status**: Active worker count and health

## 🧪 Testing

### Manual Testing with curl

1. **Create an order:**
   ```bash
   curl -X POST http://localhost:8080/orders \
     -H "Content-Type: application/json" \
     -d '{
       "amount": 150.00,
       "items": ["laptop", "mouse"],
       "customer": "test@example.com",
       "address": "456 Test St",
       "priority": 1
     }'
   ```

2. **Check statistics:**
   ```bash
   curl http://localhost:8080/stats
   ```

3. **Health check:**
   ```bash
   curl http://localhost:8080/health
   ```


