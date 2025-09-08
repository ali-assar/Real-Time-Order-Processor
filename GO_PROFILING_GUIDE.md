# Go Profiling Guide: Real-Time Order Processor

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture Deep Dive](#architecture-deep-dive)
3. [Code Walkthrough](#code-walkthrough)
4. [Go Profiling Fundamentals](#go-profiling-fundamentals)
5. [Profiling Implementation](#profiling-implementation)
6. [Practical Profiling Examples](#practical-profiling-examples)
7. [Performance Optimization Workflow](#performance-optimization-workflow)
8. [Advanced Profiling Techniques](#advanced-profiling-techniques)
9. [Troubleshooting and Best Practices](#troubleshooting-and-best-practices)

---

## Project Overview

The **Real-Time Order Processor** is a high-performance, concurrent Go application designed to process thousands of e-commerce orders per second. This project serves as an excellent learning platform for understanding Go concurrency patterns, worker pools, and most importantly, **Go profiling techniques**.

### Key Features
- **Concurrent Processing**: Worker pool architecture for high-throughput order processing
- **Priority-Based Queuing**: Orders processed based on priority levels (High/Medium/Low)
- **Real-Time Monitoring**: Live statistics and health monitoring endpoints
- **Comprehensive Profiling**: Built-in profiling endpoints for performance analysis
- **Load Testing**: Automated load generation for realistic performance testing

### Why This Project for Learning Profiling?

1. **Real Concurrency**: Multiple goroutines working simultaneously
2. **Memory Allocation**: Frequent object creation and garbage collection
3. **Channel Operations**: Blocking and non-blocking channel operations
4. **HTTP Handling**: Network I/O and request processing
5. **Business Logic**: CPU-intensive order processing

---

## Architecture Deep Dive

### System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │───▶│   Order Pool    │───▶│   Workers       │
│   (Port 8080)   │    │   (Buffer)      │    │ (Configurable)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Profiling      │    │   Results       │    │   Logging       │
│  Endpoints      │    │   Channel       │    │  & Metrics      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Component Breakdown

#### 1. **HTTP Server** (`cmd/main.go`)
- Entry point of the application
- Sets up HTTP routes and middleware
- Initializes worker pool
- Enables profiling features

#### 2. **Worker Pool** (`internal/processor/pool.go`)
- Manages concurrent order processing
- Implements priority-based queuing
- Tracks performance metrics
- Handles graceful shutdown

#### 3. **HTTP Handlers** (`internal/handler/`)
- `handler.go`: Order creation and statistics endpoints
- `profiler.go`: Profiling endpoints for performance analysis
- `router.go`: Route registration and middleware

#### 4. **Data Models** (`internal/pkg/models/model.go`)
- Order structure and validation
- Processing statistics
- Error handling

---

## Code Walkthrough

### 1. Main Application (`cmd/main.go`)

```go
func main() {
    // Enable profiling for better analysis
    runtime.SetMutexProfileFraction(1)  // Track mutex contention
    runtime.SetBlockProfileRate(1)      // Track blocking operations
    
    mux := http.NewServeMux()
    pool := processor.Start(context.Background(), 10, 100) // 10 workers, 100 buffer
    defer processor.Close(pool)
    
    // Result processing goroutine
    go func() {
        for result := range pool.Results {
            // Process results and log
        }
    }()
    
    handler.RegisterRoutes(mux, pool)
    go generateLoad(pool) // Generate test load
    
    srv := &http.Server{Addr: ":8080", Handler: mux}
    log.Fatal(srv.ListenAndServe())
}
```

**Key Profiling Points:**
- `runtime.SetMutexProfileFraction(1)`: Enables mutex profiling
- `runtime.SetBlockProfileRate(1)`: Enables block profiling
- Goroutine for result processing creates concurrency

### 2. Worker Pool (`internal/processor/pool.go`)

```go
type Pool struct {
    Orders    chan models.Order
    Results   chan models.ProcessedOrder
    Wg        sync.WaitGroup
    Ctx       context.Context
    Cancel    context.CancelFunc
    StartTime time.Time
    
    // Atomic counters for thread-safe operations
    Processed    int64
    SuccessCount int64
    ErrorCount   int64
    TotalTime    int64
    Workers      int
}
```

**Key Profiling Points:**
- **Channels**: `Orders` and `Results` channels create blocking points
- **Atomic Operations**: Thread-safe counters using `atomic` package
- **Goroutines**: Each worker runs in its own goroutine
- **Context**: Graceful shutdown mechanism

### 3. Worker Function

```go
func (p *Pool) worker(id int) {
    defer p.Wg.Done()
    for {
        select {
        case <-p.Ctx.Done():
            return
        case order, ok := <-p.Orders:
            if !ok {
                return
            }
            
            startTime := time.Now()
            processedOrder := p.processOrder(order, id, startTime)
            
            // Send result (blocking operation)
            select {
            case p.Results <- processedOrder:
            case <-p.Ctx.Done():
                return
            }
            
            // Update statistics atomically
            atomic.AddInt64(&p.Processed, 1)
            // ... more atomic operations
        }
    }
}
```

**Key Profiling Points:**
- **Channel Operations**: `<-p.Orders` and `p.Results <-` are blocking points
- **Select Statements**: Create branching in goroutine execution
- **Time Measurement**: `time.Now()` and `time.Since()` for performance tracking
- **Atomic Operations**: Thread-safe counter updates

### 4. Order Processing Logic

```go
func (p *Pool) processOrder(order models.Order, workerID int, startTime time.Time) models.ProcessedOrder {
    // Simulate processing time based on priority
    time.Sleep(time.Duration(order.Priority) * 10 * time.Millisecond)
    
    // Business logic validation
    if err := p.validateOrderForProcessing(order); err != nil {
        // Handle error
    }
    
    // Apply business rules
    if processedOrder.Success {
        processedOrder = p.applyBusinessRules(processedOrder)
    }
    
    // Calculate processing time
    processingTime := time.Since(startTime)
    processedOrder.ProcessingTime = processingTime.Milliseconds()
    
    return processedOrder
}
```

**Key Profiling Points:**
- **CPU Intensive**: Business logic processing
- **Sleep Operations**: `time.Sleep()` creates timing variations
- **Memory Allocation**: Creating `ProcessedOrder` structs
- **String Operations**: Error messages and result strings

---

## Go Profiling Fundamentals

### What is Profiling?

Profiling is the process of analyzing your program's runtime behavior to identify performance bottlenecks, memory leaks, and optimization opportunities. Go provides excellent built-in profiling tools through the `runtime/pprof` package.

### Types of Profiling in Go

#### 1. **CPU Profiling**
- **What it shows**: Where your program spends CPU time
- **When to use**: Performance optimization, bottleneck identification
- **Key metrics**: Flat time, cumulative time, call graph

#### 2. **Memory Profiling**
- **What it shows**: Memory allocation patterns and potential leaks
- **When to use**: Memory optimization, leak detection
- **Key metrics**: Allocated bytes, objects, allocation sites

#### 3. **Goroutine Profiling**
- **What it shows**: Current state of all goroutines
- **When to use**: Concurrency debugging, goroutine leaks
- **Key metrics**: Goroutine count, stack traces, states

#### 4. **Block Profiling**
- **What it shows**: Operations that block goroutines
- **When to use**: Concurrency performance issues
- **Key metrics**: Blocking operations, contention points

#### 5. **Mutex Profiling**
- **What it shows**: Mutex contention and lock contention
- **When to use**: Lock optimization, deadlock detection
- **Key metrics**: Lock contention, wait times

#### 6. **Execution Tracing**
- **What it shows**: Detailed timeline of program execution
- **When to use**: Complex concurrency debugging
- **Key metrics**: Timeline, goroutine interactions, GC events

---

## Profiling Implementation

### 1. Profiling Endpoints (`internal/handler/profiler.go`)

Our implementation provides comprehensive profiling endpoints:

```go
func RegisterProfilingRoutes(router *http.ServeMux) {
    // Standard pprof endpoints (automatically registered)
    // - /debug/pprof/ (index page)
    // - /debug/pprof/profile (CPU profile)
    // - /debug/pprof/heap (memory profile)
    // - /debug/pprof/goroutine (goroutine profile)
    // - /debug/pprof/block (block profile)
    // - /debug/pprof/mutex (mutex profile)
    // - /debug/pprof/trace (execution trace)
    
    // Custom profiling endpoints
    router.HandleFunc("/profile/cpu", CPUTraceHandler)
    router.HandleFunc("/profile/memory", MemoryTraceHandler)
    router.HandleFunc("/profile/goroutines", GoroutineTraceHandler)
    router.HandleFunc("/profile/block", BlockTraceHandler)
    router.HandleFunc("/profile/mutex", MutexTraceHandler)
    router.HandleFunc("/profile/trace", ExecutionTraceHandler)
    router.HandleFunc("/profile/gc", GCHandler)
}
```

### 2. CPU Profiling Handler

```go
func CPUTraceHandler(w http.ResponseWriter, r *http.Request) {
    duration := 30 * time.Second
    if d := r.URL.Query().Get("duration"); d != "" {
        if parsed, err := time.ParseDuration(d); err == nil {
            duration = parsed
        }
    }
    
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf("attachment; filename=cpu_profile_%d.prof", time.Now().Unix()))
    
    if err := pprof.StartCPUProfile(w); err != nil {
        http.Error(w, "failed to start CPU profile", http.StatusInternalServerError)
        return
    }
    defer pprof.StopCPUProfile()
    
    time.Sleep(duration) // Profile for specified duration
}
```

**Key Features:**
- Configurable duration via query parameter
- Automatic file download with timestamp
- Proper error handling
- Deferred cleanup

### 3. Memory Profiling Handler

```go
func MemoryTraceHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf("attachment; filename=memory_profile_%d.prof", time.Now().Unix()))
    
    runtime.GC() // Force garbage collection before profiling
    
    if err := pprof.WriteHeapProfile(w); err != nil {
        http.Error(w, "failed to write memory profile", http.StatusInternalServerError)
        return
    }
}
```

**Key Features:**
- Forces GC before profiling for accurate results
- Captures current heap state
- Thread-safe memory profiling

### 4. Goroutine Profiling Handler

```go
func GoroutineTraceHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", 
        fmt.Sprintf("attachment; filename=goroutine_profile_%d.prof", time.Now().Unix()))
    
    if err := pprof.Lookup("goroutine").WriteTo(w, 0); err != nil {
        http.Error(w, "failed to write goroutine profile", http.StatusInternalServerError)
        return
    }
}
```

**Key Features:**
- Captures current goroutine state
- Shows goroutine stack traces
- Useful for debugging concurrency issues

---

## Practical Profiling Examples

### 1. Setting Up the Environment

```bash
# Start the application
go run cmd/main.go

# In another terminal, generate load
go run tools/load_test.go

# In a third terminal, capture profiles
curl http://localhost:8080/profile/cpu -o cpu_profile.prof
curl http://localhost:8080/profile/memory -o memory_profile.prof
curl http://localhost:8080/profile/goroutines -o goroutine_profile.prof
```

### 2. CPU Profiling Analysis

```bash
# Interactive analysis
go tool pprof cpu_profile.prof

# Inside pprof interactive mode:
(pprof) top10
(pprof) list processOrder
(pprof) web
(pprof) focus processOrder
(pprof) top10
```

**Example Output:**
```
Showing nodes accounting for 2.50s, 50.0% of 5.00s total
      flat  flat%   sum%        cum   cum%
     1.20s  24.0%  24.0%      1.20s  24.0%  time.Sleep
     0.80s  16.0%  40.0%      0.80s  16.0%  processOrder
     0.30s   6.0%  46.0%      0.30s   6.0%  runtime.mallocgc
     0.20s   4.0%  50.0%      0.20s   4.0%  applyBusinessRules
```

### 3. Memory Profiling Analysis

```bash
# Interactive analysis
go tool pprof memory_profile.prof

# Inside pprof interactive mode:
(pprof) top10
(pprof) list processOrder
(pprof) web
(pprof) alloc_space
(pprof) top10
```

**Example Output:**
```
Showing nodes accounting for 512.00MB, 50.0% of 1.00GB total
      flat  flat%   sum%        cum   cum%
  256.00MB  25.0%  25.0%   256.00MB  25.0%  processOrder
  128.00MB  12.5%  37.5%   128.00MB  12.5%  createTestOrder
   64.00MB   6.3%  43.8%    64.00MB   6.3%  applyBusinessRules
```

### 4. Goroutine Profiling Analysis

```bash
# Interactive analysis
go tool pprof goroutine_profile.prof

# Inside pprof interactive mode:
(pprof) top10
(pprof) list worker
(pprof) web
```

**Example Output:**
```
Showing nodes accounting for 15, 100% of 15 total
      flat  flat%
        10  66.7%  runtime.gopark
         3  20.0%  worker
         2  13.3%  runtime.selectgo
```

### 5. Web Interface Analysis

```bash
# Start web interface
go tool pprof -http=:8081 cpu_profile.prof

# Open http://localhost:8081 in browser
# Navigate through different views:
# - Graph view: Call graph visualization
# - Flame graph: Flame graph visualization
# - Top: Top functions by resource usage
# - Source: Source code with profiling data
```

---

## Performance Optimization Workflow

### 1. **Baseline Measurement**

```bash
# Start application with load
go run cmd/main.go &
go run tools/load_test.go &

# Capture baseline profiles
curl http://localhost:8080/profile/cpu -o baseline_cpu.prof
curl http://localhost:8080/profile/memory -o baseline_memory.prof

# Analyze baseline
go tool pprof baseline_cpu.prof
(pprof) top10
(pprof) web
```

### 2. **Identify Bottlenecks**

Common bottlenecks in our order processor:

#### CPU Bottlenecks:
- **`time.Sleep()`**: Artificial delay in processing
- **`processOrder()`**: Business logic processing
- **`applyBusinessRules()`**: Rule application logic

#### Memory Bottlenecks:
- **Order creation**: Frequent struct allocation
- **String operations**: Error messages and results
- **Channel buffers**: Unbounded growth

#### Concurrency Bottlenecks:
- **Channel blocking**: Full channels causing goroutine blocking
- **Mutex contention**: Shared resource access
- **Context cancellation**: Graceful shutdown delays

### 3. **Optimize Based on Profiling Data**

#### Example Optimization 1: Reduce Memory Allocation

**Before:**
```go
func (p *Pool) processOrder(order models.Order, workerID int, startTime time.Time) models.ProcessedOrder {
    processedOrder := models.ProcessedOrder{
        Order:       order,
        ProcessedAt: time.Now(),
        WorkerID:    workerID,
        Success:     true,
        Result:      "Order processed successfully", // String allocation
    }
    // ... more allocations
}
```

**After:**
```go
// Pre-allocate common strings
var (
    successResult = "Order processed successfully"
    errorResult   = "Order processing failed"
)

func (p *Pool) processOrder(order models.Order, workerID int, startTime time.Time) models.ProcessedOrder {
    processedOrder := models.ProcessedOrder{
        Order:       order,
        ProcessedAt: time.Now(),
        WorkerID:    workerID,
        Success:     true,
        Result:      successResult, // Reuse pre-allocated string
    }
    // ... reuse more pre-allocated strings
}
```

#### Example Optimization 2: Optimize Channel Operations

**Before:**
```go
// Blocking channel operation
case p.Results <- processedOrder:
```

**After:**
```go
// Non-blocking with timeout
select {
case p.Results <- processedOrder:
case <-time.After(100 * time.Millisecond):
    // Handle timeout
case <-p.Ctx.Done():
    return
}
```

### 4. **Verify Improvements**

```bash
# Capture optimized profiles
curl http://localhost:8080/profile/cpu -o optimized_cpu.prof
curl http://localhost:8080/profile/memory -o optimized_memory.prof

# Compare with baseline
go tool pprof -base baseline_cpu.prof optimized_cpu.prof
(pprof) top10
(pprof) web
```

---

## Advanced Profiling Techniques

### 1. **Continuous Profiling**

Set up automated profiling collection:

```bash
#!/bin/bash
# continuous_profiling.sh

while true; do
    timestamp=$(date +%s)
    
    # Collect profiles every minute
    curl -s http://localhost:8080/profile/cpu -o "profiles/cpu_${timestamp}.prof"
    curl -s http://localhost:8080/profile/memory -o "profiles/memory_${timestamp}.prof"
    curl -s http://localhost:8080/profile/goroutines -o "profiles/goroutine_${timestamp}.prof"
    
    sleep 60
done
```

### 2. **Profile Comparison**

```bash
# Compare two profiles
go tool pprof -base old_profile.prof new_profile.prof

# Compare multiple profiles
go tool pprof -base baseline.prof -base optimized.prof final.prof
```

### 3. **Custom Profiling Metrics**

Add custom metrics to your application:

```go
// Add to Pool struct
type Pool struct {
    // ... existing fields
    
    // Custom profiling metrics
    ChannelBlockTime int64
    ProcessingTime   int64
    QueueWaitTime    int64
}

// Add to worker function
func (p *Pool) worker(id int) {
    for {
        select {
        case order, ok := <-p.Orders:
            start := time.Now()
            
            // Measure channel wait time
            atomic.AddInt64(&p.QueueWaitTime, time.Since(start).Nanoseconds())
            
            // Process order
            processedOrder := p.processOrder(order, id, start)
            
            // Measure processing time
            atomic.AddInt64(&p.ProcessingTime, processedOrder.ProcessingTime)
            
            // ... rest of processing
        }
    }
}
```

### 4. **Profiling in Production**

```go
// Add to main.go
func main() {
    // Enable profiling in production
    if os.Getenv("ENABLE_PROFILING") == "true" {
        go func() {
            log.Println("Profiling server starting on :6060")
            log.Fatal(http.ListenAndServe(":6060", nil))
        }()
    }
    
    // ... rest of main
}
```

### 5. **Memory Leak Detection**

```go
// Add to Pool struct
type Pool struct {
    // ... existing fields
    
    // Memory leak detection
    lastGC     time.Time
    gcCount    int64
    leakCount  int64
}

// Add to worker function
func (p *Pool) worker(id int) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Check for memory leaks
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            if m.NumGC > p.gcCount {
                p.gcCount = m.NumGC
                p.lastGC = time.Now()
            } else if time.Since(p.lastGC) > 5*time.Minute {
                // Potential memory leak detected
                atomic.AddInt64(&p.leakCount, 1)
                log.Printf("Potential memory leak detected: %d leaks", p.leakCount)
            }
            
        case order, ok := <-p.Orders:
            // ... process order
        }
    }
}
```

---

## Troubleshooting and Best Practices

### Common Profiling Issues

#### 1. **"No samples found"**
- **Cause**: Application not under load
- **Solution**: Generate load before profiling
- **Prevention**: Always profile under realistic conditions

#### 2. **"Failed to start CPU profiling"**
- **Cause**: Another profile already running
- **Solution**: Wait for current profile to complete
- **Prevention**: Implement profile locking

#### 3. **"Graphviz not found"**
- **Cause**: Missing graphviz installation
- **Solution**: Install graphviz or use `png` command
- **Prevention**: Document system requirements

### Best Practices

#### 1. **Profile Under Realistic Load**
```bash
# Good: Profile under load
go run cmd/main.go &
go run tools/load_test.go &
curl http://localhost:8080/profile/cpu -o cpu.prof

# Bad: Profile without load
go run cmd/main.go &
curl http://localhost:8080/profile/cpu -o cpu.prof
```

#### 2. **Profile Multiple Times**
```bash
# Capture multiple profiles for consistency
for i in {1..5}; do
    curl http://localhost:8080/profile/cpu -o "cpu_${i}.prof"
    sleep 10
done
```

#### 3. **Use Multiple Profile Types**
```bash
# Comprehensive profiling
curl http://localhost:8080/profile/cpu -o cpu.prof
curl http://localhost:8080/profile/memory -o memory.prof
curl http://localhost:8080/profile/goroutines -o goroutine.prof
curl http://localhost:8080/profile/block -o block.prof
```

#### 4. **Profile in Production**
```go
// Enable production profiling
if os.Getenv("ENV") == "production" {
    go func() {
        log.Println("Production profiling enabled")
        http.ListenAndServe(":6060", nil)
    }()
}
```

#### 5. **Automate Profiling**
```yaml
# docker-compose.yml
version: '3'
services:
  app:
    build: .
    ports:
      - "8080:8080"
      - "6060:6060"  # Profiling port
    environment:
      - ENABLE_PROFILING=true
```

### Performance Optimization Checklist

- [ ] **CPU Profiling**: Identify hot functions
- [ ] **Memory Profiling**: Find allocation hotspots
- [ ] **Goroutine Profiling**: Check for leaks
- [ ] **Block Profiling**: Identify blocking operations
- [ ] **Mutex Profiling**: Check for contention
- [ ] **Execution Tracing**: Understand concurrency patterns
- [ ] **Load Testing**: Verify under realistic conditions
- [ ] **Continuous Monitoring**: Set up automated profiling
- [ ] **Baseline Comparison**: Measure before/after improvements
- [ ] **Production Profiling**: Monitor in production

---

## Conclusion

This Real-Time Order Processor project provides an excellent foundation for learning Go profiling. The combination of:

- **Concurrent processing** with worker pools
- **Channel-based communication** between goroutines
- **HTTP request handling** with network I/O
- **Business logic processing** with CPU-intensive operations
- **Memory allocation** patterns from order processing

Creates realistic scenarios where profiling techniques can be applied and learned effectively.

### Key Takeaways

1. **Profiling is essential** for performance optimization
2. **Multiple profile types** provide different insights
3. **Profile under realistic load** for accurate results
4. **Continuous profiling** helps catch performance regressions
5. **Production profiling** enables real-world optimization

### Next Steps

1. **Experiment** with different profiling scenarios
2. **Optimize** based on profiling data
3. **Implement** continuous profiling in production
4. **Learn** advanced profiling techniques
5. **Apply** profiling to your own projects

Remember: **Profiling is not a one-time activity** - it's an ongoing process that should be integrated into your development workflow for optimal performance.

---

## Quick Reference

### Profiling Commands
```bash
# Start application
go run cmd/main.go

# Generate load
go run tools/load_test.go

# Capture profiles
curl http://localhost:8080/profile/cpu -o cpu.prof
curl http://localhost:8080/profile/memory -o memory.prof
curl http://localhost:8080/profile/goroutines -o goroutine.prof

# Analyze profiles
go tool pprof cpu.prof
go tool pprof -http=:8081 cpu.prof

# Compare profiles
go tool pprof -base old.prof new.prof
```

### Useful pprof Commands
```
top10          # Top 10 functions
list function  # Show source code
web            # Generate call graph
png            # Generate PNG
focus func     # Focus on function
ignore func    # Ignore function
```

### Profiling Endpoints
- `/debug/pprof/` - Main profiling page
- `/profile/cpu` - CPU profile
- `/profile/memory` - Memory profile
- `/profile/goroutines` - Goroutine profile
- `/profile/block` - Block profile
- `/profile/mutex` - Mutex profile
- `/profile/trace` - Execution trace
- `/profile/gc` - GC statistics

Happy profiling! 🚀
