package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/processor"
)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request, pool *processor.Pool) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var o models.Order
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Set default values before validation
	o.SetDefaultValues()

	// Generate ID if not provided
	if o.ID == "" {
		o.ID = generateID()
	}

	if err := o.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set creation time
	o.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(o)

	// Send to processing pool
	select {
	case pool.Orders <- o:
		// Order queued successfully
	default:
		// Queue is full
		http.Error(w, "service temporarily unavailable", http.StatusServiceUnavailable)
		return
	}
}

func generateID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// GetStatsHandler returns processing statistics
func GetStatsHandler(w http.ResponseWriter, r *http.Request, pool *processor.Pool) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	stats := pool.Stats()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stats)
}

// HealthCheckHandler returns the health status of the service
func HealthCheckHandler(w http.ResponseWriter, r *http.Request, pool *processor.Pool) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"pool": map[string]interface{}{
			"healthy":      pool.IsHealthy(),
			"queue_length": pool.GetQueueLength(),
			"workers":      pool.Workers,
		},
	}

	if !pool.IsHealthy() {
		health["status"] = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(health)
}
