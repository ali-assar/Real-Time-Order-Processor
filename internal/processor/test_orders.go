package processor

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
)

// CreateTestOrder generates a test order for load testing
func CreateTestOrder(id int) models.Order {
	rand.Seed(time.Now().UnixNano())

	items := []string{"laptop", "mouse", "keyboard", "monitor", "headphones", "webcam", "speaker", "tablet"}
	itemCount := rand.Intn(5) + 1 

	orderItems := make([]string, itemCount)
	for i := 0; i < itemCount; i++ {
		orderItems[i] = items[rand.Intn(len(items))]
	}

	priorities := []int{1, 2, 3} // High, Medium, Low
	priority := priorities[rand.Intn(len(priorities))]

	// Create some high-value orders for priority processing
	amount := float64(rand.Intn(2000) + 10) // $10-$2010
	if rand.Float64() < 0.1 {               // 10% chance of high-value order
		amount = float64(rand.Intn(5000) + 1000) // $1000-$6000
		priority = 1                             // High priority for high-value orders
	}

	return models.Order{
		ID:       fmt.Sprintf("test_order_%d", id),
		Amount:   amount,
		Items:    orderItems,
		Customer: fmt.Sprintf("customer%d@example.com", rand.Intn(1000)),
		Status:   "pending",
		Address:  fmt.Sprintf("%d Test Street, City %d", rand.Intn(1000), rand.Intn(100)),
		Priority: priority,
		Notes:    fmt.Sprintf("Test order %d", id),
	}
}
