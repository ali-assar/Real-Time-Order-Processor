package processor

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
)

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
	TotalTime    int64 // total processing time in milliseconds

	Workers int
}

func Start(ctx context.Context, workers, buf int) *Pool {
	ctx, cancel := context.WithCancel(ctx)
	pool := &Pool{
		Orders:    make(chan models.Order, buf),
		Results:   make(chan models.ProcessedOrder, buf),
		Workers:   workers,
		Ctx:       ctx,
		Cancel:    cancel,
		StartTime: time.Now(),
	}

	for i := 0; i < workers; i++ {
		pool.Wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func Close(pool *Pool) {
	pool.Cancel()
	pool.Wg.Wait()
	close(pool.Orders)
	close(pool.Results)
}

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

			// Send result to results channel
			select {
			case p.Results <- processedOrder:
			case <-p.Ctx.Done():
				return
			}

			// Update statistics
			atomic.AddInt64(&p.Processed, 1)
			if processedOrder.Success {
				atomic.AddInt64(&p.SuccessCount, 1)
			} else {
				atomic.AddInt64(&p.ErrorCount, 1)
			}
			atomic.AddInt64(&p.TotalTime, processedOrder.ProcessingTime)
		}
	}
}

func (p *Pool) processOrder(order models.Order, workerID int, startTime time.Time) models.ProcessedOrder {
	processedOrder := models.ProcessedOrder{
		Order:       order,
		ProcessedAt: time.Now(),
		WorkerID:    workerID,
		Success:     true,
		Result:      "Order processed successfully",
	}

	// Simulate order processing logic
	time.Sleep(time.Duration(order.Priority) * 10 * time.Millisecond) // Priority-based processing time

	// Business logic validation and processing
	if err := p.validateOrderForProcessing(order); err != nil {
		processedOrder.Success = false
		processedOrder.Error = err.Error()
		processedOrder.Result = "Order processing failed"
	}

	// Simulate additional processing steps
	if processedOrder.Success {
		processedOrder = p.applyBusinessRules(processedOrder)
	}

	// Calculate processing time
	processingTime := time.Since(startTime)
	processedOrder.ProcessingTime = processingTime.Milliseconds()

	return processedOrder
}

func (p *Pool) validateOrderForProcessing(order models.Order) error {
	// Additional business validation
	if order.Amount > 10000 {
		return &models.ValidationError{Message: "order amount exceeds limit"}
	}

	if len(order.Items) > 50 {
		return &models.ValidationError{Message: "too many items in order"}
	}

	return nil
}

func (p *Pool) applyBusinessRules(processedOrder models.ProcessedOrder) models.ProcessedOrder {
	order := &processedOrder.Order

	// Apply business rules based on order characteristics
	switch {
	case order.Amount > 1000:
		order.Status = "priority_processing"
		processedOrder.Result = "Order marked for priority processing"
	case order.Priority == 1:
		order.Status = "expedited"
		processedOrder.Result = "Order expedited due to high priority"
	default:
		order.Status = "processing"
		processedOrder.Result = "Order processing completed"
	}

	return processedOrder
}

func (p *Pool) Stats() models.ProcessingStats {
	processed := atomic.LoadInt64(&p.Processed)
	success := atomic.LoadInt64(&p.SuccessCount)
	error := atomic.LoadInt64(&p.ErrorCount)
	totalTime := atomic.LoadInt64(&p.TotalTime)

	var avgTime float64
	if processed > 0 {
		avgTime = float64(totalTime) / float64(processed)
	}

	uptime := int64(time.Since(p.StartTime).Seconds())

	return models.ProcessingStats{
		TotalProcessed:     int(processed),
		SuccessCount:       int(success),
		ErrorCount:         int(error),
		AverageProcessTime: avgTime,
		ActiveWorkers:      p.Workers,
		QueueLength:        len(p.Orders),
		Uptime:             uptime,
	}
}

// GetQueueLength returns the current number of orders in the queue
func (p *Pool) GetQueueLength() int {
	return len(p.Orders)
}

// IsHealthy checks if the pool is in a healthy state
func (p *Pool) IsHealthy() bool {
	return p.Ctx.Err() == nil && len(p.Orders) < cap(p.Orders)
}
