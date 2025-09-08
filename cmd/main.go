package main

import (
	"context"
	"log"
	"net/http"
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("API listening on %s", srv.Addr)
	log.Printf("Profiling available at http://localhost:8080/debug/pprof/")
	log.Fatal(srv.ListenAndServe())
}
