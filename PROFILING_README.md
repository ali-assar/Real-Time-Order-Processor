# Go Profiling with Real-Time Order Processor

This project demonstrates comprehensive Go profiling techniques using a realistic, high-performance order processing system.

## 🚀 Quick Start

### Option 1: Automated Script (Recommended)

**Linux/macOS:**
```bash
chmod +x scripts/start_profiling.sh
./scripts/start_profiling.sh
```

**Windows:**
```cmd
scripts\start_profiling.bat
```

### Option 2: Manual Setup

1. **Start the application:**
   ```bash
   go run cmd/main.go
   ```

2. **Generate load (in another terminal):**
   ```bash
   go run tools/load_test.go
   ```

3. **Capture profiles:**
   ```bash
   # CPU profile (30 seconds)
   curl http://localhost:8080/profile/cpu -o cpu_profile.prof
   
   # Memory profile
   curl http://localhost:8080/profile/memory -o memory_profile.prof
   
   # Goroutine profile
   curl http://localhost:8080/profile/goroutines -o goroutine_profile.prof
   ```

4. **Analyze profiles:**
   ```bash
   # Interactive analysis
   go tool pprof cpu_profile.prof
   
   # Web interface
   go tool pprof -http=:8081 cpu_profile.prof
   ```

## 📊 Available Profiling Endpoints

| Endpoint | Description | Duration |
|----------|-------------|----------|
| `/debug/pprof/` | Main profiling page | - |
| `/profile/cpu` | CPU profile | 30s (configurable) |
| `/profile/memory` | Memory profile | Instant |
| `/profile/goroutines` | Goroutine profile | Instant |
| `/profile/block` | Block profile | Instant |
| `/profile/mutex` | Mutex profile | Instant |
| `/profile/trace` | Execution trace | 5s (configurable) |
| `/profile/gc` | GC statistics | Instant |

## 🔍 Profiling Analysis

### Interactive Commands

```bash
go tool pprof cpu_profile.prof
```

**Inside pprof:**
```
(pprof) top10          # Top 10 functions by resource usage
(pprof) list processOrder  # Show source code with profiling data
(pprof) web            # Generate call graph (requires graphviz)
(pprof) png            # Generate PNG visualization
(pprof) focus processOrder  # Focus on specific function
(pprof) ignore runtime  # Ignore runtime functions
```

### Web Interface

```bash
go tool pprof -http=:8081 cpu_profile.prof
```

Open http://localhost:8081 in your browser for:
- **Graph view**: Call graph visualization
- **Flame graph**: Flame graph visualization  
- **Top view**: Top functions by resource usage
- **Source view**: Source code with profiling data

## 📈 Performance Scenarios

The application includes three load testing scenarios:

1. **Normal Load**: 10 orders/sec for 30 seconds
2. **High Load**: 50 orders/sec for 20 seconds  
3. **Burst Load**: 100 orders/sec for 10 seconds

## 🛠️ Troubleshooting

### Common Issues

**"No samples found"**
- Ensure the application is under load
- Wait for load generation to start

**"Failed to start CPU profiling"**
- Check if another profile is running
- Verify the application is responding

**"Graphviz not found"**
- Install graphviz: `brew install graphviz` (macOS) or `apt-get install graphviz` (Ubuntu)
- Use `png` command instead of `web`

### System Requirements

- Go 1.24.4 or higher
- curl (for profile capture)
- graphviz (for web visualizations, optional)

## 📚 Learning Resources

- **Comprehensive Guide**: See `GO_PROFILING_GUIDE.md` for detailed explanations
- **Go pprof Documentation**: https://golang.org/pkg/runtime/pprof/
- **Go Profiling Blog**: https://golang.org/doc/diagnostics.html

## 🎯 What You'll Learn

- **CPU Profiling**: Identify performance bottlenecks
- **Memory Profiling**: Find memory leaks and allocation hotspots
- **Goroutine Profiling**: Debug concurrency issues
- **Block Profiling**: Identify blocking operations
- **Mutex Profiling**: Find lock contention
- **Execution Tracing**: Understand complex concurrency patterns

## 🔧 Customization

### Modify Load Patterns

Edit `cmd/load_test.go` to create custom load scenarios:

```go
// Custom scenario
runLoadTest(baseURL, 25, 60*time.Second, "custom")
```

### Adjust Worker Pool

Modify `cmd/main.go` to change worker configuration:

```go
pool := processor.Start(context.Background(), 20, 200) // 20 workers, 200 buffer
```

### Add Custom Metrics

Extend the `Pool` struct in `internal/processor/pool.go` with custom profiling metrics.

## 📊 Example Analysis Session

```bash
# 1. Start application and load
go run cmd/main.go &
go run cmd/load_test.go &

# 2. Capture profiles
curl http://localhost:8080/profile/cpu -o cpu.prof
curl http://localhost:8080/profile/memory -o memory.prof

# 3. Analyze CPU profile
go tool pprof cpu.prof
(pprof) top10
(pprof) list processOrder
(pprof) web

# 4. Analyze memory profile
go tool pprof memory.prof
(pprof) top10
(pprof) alloc_space
(pprof) web

# 5. Compare profiles
go tool pprof -base old_cpu.prof cpu.prof
```

