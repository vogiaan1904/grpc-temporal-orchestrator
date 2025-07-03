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

const PrePaymentOrderTaskQueue = "PRE_PAYMENT_ORDER_TASK_QUEUE"
const PostPaymentOrderTaskQueue = "POST_PAYMENT_ORDER_TASK_QUEUE"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	c, err := tClient.Dial(tClient.Options{
		HostPort:  cfg.Temporal.HostPort,
		Namespace: cfg.Temporal.Namespace,
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	grpcClis, cleanup, err := pkgGrpc.InitGrpcClients(pkgGrpc.GrpcAddresses{
		ProductAddr: cfg.Grpc.ProductSvcAddr,
		OrderAddr:   cfg.Grpc.OrderSvcAddr,
		PaymentAddr: cfg.Grpc.PaymentSvcAddr,
	})
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer cleanup()

	prePaymentWorker := worker.New(c, PrePaymentOrderTaskQueue, worker.Options{})
	prePaymentWorker.RegisterWorkflow(oWF.ProcessPrePaymentOrder)

	postPaymentWorker := worker.New(c, PostPaymentOrderTaskQueue, worker.Options{})
	postPaymentWorker.RegisterWorkflow(oWF.ProcessPostPaymentOrder)

	oActs := &activities.OrderActivities{
		Client: grpcClis.Order,
	}
	paymentActs := &activities.PaymentActivities{
		Client: grpcClis.Payment,
	}
	prodActs := &activities.ProductActivities{
		Client: grpcClis.Product,
	}

	prePaymentWorker.RegisterActivity(oActs)
	prePaymentWorker.RegisterActivity(paymentActs)
	prePaymentWorker.RegisterActivity(prodActs)

	postPaymentWorker.RegisterActivity(oActs)
	postPaymentWorker.RegisterActivity(paymentActs)
	postPaymentWorker.RegisterActivity(prodActs)

	err = prePaymentWorker.Start()
	if err != nil {
		log.Fatalf("Failed to start pre-payment worker: %v", err)
	}
	log.Println("Pre-payment worker started")

	err = postPaymentWorker.Start()
	if err != nil {
		log.Fatalf("Failed to start post-payment worker: %v", err)
	}
	log.Println("Post-payment worker started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down worker...")
}
