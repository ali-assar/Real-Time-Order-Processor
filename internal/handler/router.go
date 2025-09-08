package handler

import (
	"net/http"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func RegisterRoutes(router *http.ServeMux, pool *processor.Pool) {
	// Order management
	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateOrderHandler(w, r, pool)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Statistics and monitoring
	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		GetStatsHandler(w, r, pool)
	})

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		HealthCheckHandler(w, r, pool)
	})
}
