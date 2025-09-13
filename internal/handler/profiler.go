package handler

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func RegisterProfilingRoutes(router *http.ServeMux) {
	// Standard pprof endpoints (automatically registered by importing net/http/pprof)
	// These are available at:
	// - /debug/pprof/ (index page)
	// - /debug/pprof/profile (CPU profile)
	// - /debug/pprof/heap (memory profile)
	// - /debug/pprof/goroutine (goroutine profile)
	// - /debug/pprof/block (block profile)
	// - /debug/pprof/mutex (mutex profile)
	// - /debug/pprof/trace (execution trace)

	router.HandleFunc("/profile/cpu", CPUTraceHandler)
	router.HandleFunc("/profile/memory", MemoryTraceHandler)
	router.HandleFunc("/profile/goroutine", GoroutineTraceHandler)
	router.HandleFunc("/profile/block", BlockTraceHandler)
	router.HandleFunc("/profile/mutex", MutexTraceHandler)
	router.HandleFunc("/profile/trace", ExecutionTraceHandler)
	router.HandleFunc("/profile/gc", GCHandler)

}

func CPUTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	duration := 30 * time.Second
	if d := r.URL.Query().Get("duration"); d != "" {
		dur, err := time.ParseDuration(d)
		if err != nil {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		duration = dur
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=cpu_profile_%d.prof", time.Now().Unix()))

	if err := pprof.StartCPUProfile(w); err != nil {
		http.Error(w, "failed to start CPU profile", http.StatusInternalServerError)
		return
	}
	defer pprof.StopCPUProfile()
	time.Sleep(duration)
}

func MemoryTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=memory_profile_%d.prof", time.Now().Unix()))

	runtime.GC() // Force garbage collection before capturing memory profile

	if err := pprof.WriteHeapProfile(w); err != nil {
		http.Error(w, "failed to write memory profile", http.StatusInternalServerError)
		return
	}
}

func GoroutineTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=goroutine_profile_%d.prof", time.Now().Unix()))

	if err := pprof.Lookup("goroutine").WriteTo(w, 0); err != nil {
		http.Error(w, "failed to write goroutine profile", http.StatusInternalServerError)
		return
	}
}

// BlockTraceHandler captures block profile
func BlockTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=block_profile_%d.prof", time.Now().Unix()))

	if err := pprof.Lookup("block").WriteTo(w, 0); err != nil {
		http.Error(w, "failed to write block profile", http.StatusInternalServerError)
		return
	}
}

func MutexTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=mutex_profile_%d.prof", time.Now().Unix()))

	if err := pprof.Lookup("mutex").WriteTo(w, 0); err != nil {
		http.Error(w, "failed to write mutex profile", http.StatusInternalServerError)
		return
	}
}

// ExecutionTraceHandler starts execution tracing
func ExecutionTraceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	duration := 5 * time.Second
	if d := r.URL.Query().Get("duration"); d != "" {
		if parsed, err := time.ParseDuration(d); err == nil {
			duration = parsed
		}
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=trace_%d.trace", time.Now().Unix()))
	if err := trace.Start(w); err != nil {
		http.Error(w, "Failed to start trace", http.StatusInternalServerError)
		return
	}

	time.Sleep(duration)
	trace.Stop()

}

// GCHandler triggers garbage collection and shows GC stats
func GCHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get GC stats before
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// Force garbage collection
	runtime.GC()

	// Get GC stats after
	runtime.ReadMemStats(&m2)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"before_gc": {
			"alloc_mb": %.2f,
			"total_alloc_mb": %.2f,
			"sys_mb": %.2f,
			"num_gc": %d
		},
		"after_gc": {
			"alloc_mb": %.2f,
			"total_alloc_mb": %.2f,
			"sys_mb": %.2f,
			"num_gc": %d
		},
		"gc_duration_ms": %.2f
	}`,
		float64(m1.Alloc)/1024/1024,
		float64(m1.TotalAlloc)/1024/1024,
		float64(m1.Sys)/1024/1024,
		m1.NumGC,
		float64(m2.Alloc)/1024/1024,
		float64(m2.TotalAlloc)/1024/1024,
		float64(m2.Sys)/1024/1024,
		m2.NumGC,
		float64(m2.PauseTotalNs)/1000000,
	)
}
