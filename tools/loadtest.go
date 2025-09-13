package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ali-assar/Real-Time-Order-Processor.git/internal/pkg/models"
)

func main() {
	baseURL := "http://localhost:8080"

	// Test different scenarios
	fmt.Println("Starting load test...")

	// Scenario 1: Normal load
	fmt.Println("Scenario 1: Normal load (10 orders/sec for 30 seconds)")
	runLoadTest(baseURL, 10, 30*time.Second, "normal")

	time.Sleep(5 * time.Second)

	// Scenario 2: High load
	fmt.Println("Scenario 2: High load (50 orders/sec for 20 seconds)")
	runLoadTest(baseURL, 50, 20*time.Second, "high")

	time.Sleep(5 * time.Second)

	// Scenario 3: Burst load
	fmt.Println("Scenario 3: Burst load (100 orders/sec for 10 seconds)")
	runLoadTest(baseURL, 100, 10*time.Second, "burst")

	fmt.Println("Load test completed!")
}

func runLoadTest(baseURL string, ordersPerSecond int, duration time.Duration, scenario string) {
	var wg sync.WaitGroup
	interval := time.Second / time.Duration(ordersPerSecond)

	start := time.Now()
	orderCount := 0

	for time.Since(start) < duration {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendOrder(baseURL, orderCount, scenario)
		}()

		orderCount++
		time.Sleep(interval)
	}

	wg.Wait()
	fmt.Printf("Sent %d orders in %s scenario\n", orderCount, scenario)
}

func sendOrder(baseURL string, id int, scenario string) {
	order := createTestOrder(id, scenario)

	jsonData, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshaling order: %v", err)
		return
	}

	resp, err := http.Post(baseURL+"/orders", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending order: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Order %d failed with status: %d", id, resp.StatusCode)
	}
}

func createTestOrder(id int, scenario string) models.Order {
	// Create different order patterns based on scenario
	switch scenario {
	case "normal":
		return models.Order{
			ID:       fmt.Sprintf("normal_%d", id),
			Amount:   float64(50 + id%200), // $50-$250
			Items:    []string{"item1", "item2"},
			Customer: fmt.Sprintf("customer%d@example.com", id),
			Status:   "pending",
			Address:  fmt.Sprintf("%d Main St", id),
			Priority: 2, // Medium priority
			Notes:    "Normal order",
		}
	case "high":
		return models.Order{
			ID:       fmt.Sprintf("high_%d", id),
			Amount:   float64(500 + id%1000), // $500-$1500
			Items:    []string{"expensive_item1", "expensive_item2"},
			Customer: fmt.Sprintf("vip_customer%d@example.com", id),
			Status:   "pending",
			Address:  fmt.Sprintf("%d VIP Street", id),
			Priority: 1, // High priority
			Notes:    "High value order",
		}
	case "burst":
		return models.Order{
			ID:       fmt.Sprintf("burst_%d", id),
			Amount:   float64(10 + id%100), // $10-$110
			Items:    []string{"quick_item"},
			Customer: fmt.Sprintf("burst_customer%d@example.com", id),
			Status:   "pending",
			Address:  fmt.Sprintf("%d Quick St", id),
			Priority: 3, // Low priority
			Notes:    "Burst order",
		}
	default:
		return models.Order{
			ID:       fmt.Sprintf("test_%d", id),
			Amount:   100.0,
			Items:    []string{"test_item"},
			Customer: "test@example.com",
			Status:   "pending",
			Address:  "Test Address",
			Priority: 2,
			Notes:    "Test order",
		}
	}
}
