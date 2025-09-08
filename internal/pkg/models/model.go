package models

import (
	"errors"
	"time"
)

type Order struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Items     []string  `json:"items"`
	Customer  string    `json:"customer"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Address   string    `json:"address"`
	Notes     string    `json:"notes,omitempty"`
	Priority  int       `json:"priority,omitempty"` // 1=high, 2=medium, 3=low
}

type ProcessedOrder struct {
	Order          Order     `json:"order"`
	ProcessedAt    time.Time `json:"processed_at"`
	ProcessingTime int64     `json:"processing_time_ms"`
	WorkerID       int       `json:"worker_id"`
	Success        bool      `json:"success"`
	Error          string    `json:"error,omitempty"`
	Result         string    `json:"result,omitempty"`
}

type ProcessingStats struct {
	TotalProcessed     int     `json:"total_processed"`
	SuccessCount       int     `json:"success_count"`
	ErrorCount         int     `json:"error_count"`
	AverageProcessTime float64 `json:"average_process_time_ms"`
	ActiveWorkers      int     `json:"active_workers"`
	QueueLength        int     `json:"queue_length"`
	Uptime             int64   `json:"uptime_seconds"`
}

var validStatuses = map[string]bool{
	"pending":   true,
	"paid":      true,
	"shipped":   true,
	"delivered": true,
	"cancelled": true,
}

var validPriorities = map[int]bool{
	1: true, // high
	2: true, // medium
	3: true, // low
}

// Validate checks whether the order has all required fields with acceptable values.
func (o *Order) Validate() error {
	if o.ID == "" {
		return errors.New("id is required")
	}
	if o.Customer == "" {
		return errors.New("customer is required")
	}
	if o.Status == "" {
		return errors.New("status is required")
	}
	if !validStatuses[o.Status] {
		return errors.New("invalid status")
	}
	if o.Address == "" {
		return errors.New("address is required")
	}
	if o.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if len(o.Items) == 0 {
		return errors.New("items must not be empty")
	}
	if o.Priority != 0 && !validPriorities[o.Priority] {
		return errors.New("invalid priority (must be 1, 2, or 3)")
	}
	return nil
}

// SetDefaultValues sets default values for optional fields
func (o *Order) SetDefaultValues() {
	if o.Priority == 0 {
		o.Priority = 2 // default to medium priority
	}
	if o.Status == "" {
		o.Status = "pending"
	}
}

type ValidationError struct {
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
