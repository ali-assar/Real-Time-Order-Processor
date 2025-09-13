package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof" // Import for side effects - registers pprof handlers
	"runtime"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/handler"
	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func main() {

	// Enable mutex profiling for better analysis
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	mux := http.NewServeMux()
	pool := processor.Start(context.Background(), 10, 100)
	defer processor.Close(pool)

	// Start result processor goroutine
	go func() {
		for result := range pool.Results {
			if result.Success {
				log.Printf("✅ Order %s processed successfully by worker %d in %dms: %s",
					result.Order.ID, result.WorkerID, result.ProcessingTime, result.Result)
			} else {
				log.Printf("❌ Order %s processing failed by worker %d: %s",
					result.Order.ID, result.WorkerID, result.Error)
			}
		}
	}()

	handler.RegisterRoutes(mux, pool)

	// Register pprof handlers with our custom mux
	// The pprof package automatically registers handlers with http.DefaultServeMux
	// We need to mount them on our custom mux
	mux.Handle("/debug/pprof/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/heap", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/goroutine", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/block", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/mutex", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("API listening on %s", srv.Addr)
	log.Printf("Profiling available at http://localhost:8080/debug/pprof/")
	log.Fatal(srv.ListenAndServe())
}
