package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vogiaan1904/order-orchestrator/config"
	"github.com/vogiaan1904/order-orchestrator/internal/activities"
	oWF "github.com/vogiaan1904/order-orchestrator/internal/workflows/order"
	pkgGrpc "github.com/vogiaan1904/order-orchestrator/pkg/grpc"
	tClient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "ORDER_PROCESSING_TASK_QUEUE"

func main() {
	log.Println("Starting Order Orchestrator Worker")

	// Create Temporal client
	c, err := tClient.Dial(tClient.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	grpcClis, cleanup, err := pkgGrpc.InitGrpcClients(pkgGrpc.GrpcAddresses{
		ProductAddr: cfg.Grpc.ProductSvcAddr,
		OrderAddr:   cfg.Grpc.OrderSvcAddr,
		PaymentAddr: cfg.Grpc.PaymentSvcAddr,
	})
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer cleanup()

	// Create worker
	w := worker.New(c, TaskQueue, worker.Options{})
	// Register workflows and activities
	w.RegisterWorkflow(oWF.ProcessPrePaymentOrder)
	w.RegisterWorkflow(oWF.ProcessPostPaymentOrder)

	oActs := &activities.OrderActivities{
		Client: grpcClis.Order,
	}
	pActs := &activities.PaymentActivities{
		Client: grpcClis.Payment,
	}

	w.RegisterActivity(oActs)
	w.RegisterActivity(pActs)

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
