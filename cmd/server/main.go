package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vogiaan1904/order-orchestrator/internal/activities"
	"github.com/vogiaan1904/order-orchestrator/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "ORDER_PROCESSING_TASK_QUEUE"

func main() {
	log.Println("Starting Order Orchestrator Worker")

	// Create Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// Configure service addresses from environment or use defaults
	orderSvcAddr := getEnv("ORDER_SERVICE_ADDRESS", "localhost:50052")
	paymentSvcAddr := getEnv("PAYMENT_SERVICE_ADDRESS", "localhost:50055")

	// Create worker
	w := worker.New(c, TaskQueue, worker.Options{})

	// Register workflows and activities
	w.RegisterWorkflow(workflows.OrderProcessingWorkflow)

	orderActivities := &activities.OrderActivities{
		OrderSvcAddr: orderSvcAddr,
	}
	paymentActivities := &activities.PaymentActivities{
		PaymentSvcAddr: paymentSvcAddr,
	}

	w.RegisterActivity(orderActivities)
	w.RegisterActivity(paymentActivities)

	// Start worker (non-blocking)
	err = w.Start()
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
	log.Println("Worker started")

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down worker...")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
