package grpcservices

import (
	"log"

	"github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"github.com/vogiaan1904/order-orchestrator/protogen/golang/payment"
	"github.com/vogiaan1904/order-orchestrator/protogen/golang/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	Product product.ProductServiceClient
	Order   order.OrderServiceClient
	Payment payment.PaymentServiceClient
}

type GrpcAddresses struct {
	ProductAddr string
	OrderAddr   string
	PaymentAddr string
}

type cleanupFunc func()

func InitGrpcClients(addresses GrpcAddresses) (*GrpcClients, cleanupFunc, error) {
	var cleanupFuncs []cleanupFunc
	clients := &GrpcClients{}

	// Product gRPC client
	productConn, err := grpc.NewClient(
		addresses.ProductAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	clients.Product = product.NewProductServiceClient(productConn)

	// Order gRPC client
	orderConn, err := grpc.NewClient(
		addresses.OrderAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	clients.Order = order.NewOrderServiceClient(orderConn)

	// Payment gRPC client
	paymentConn, err := grpc.NewClient(
		addresses.PaymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}
	clients.Payment = payment.NewPaymentServiceClient(paymentConn)

	// Cleanup function
	cleanupFuncs = append(cleanupFuncs, func() {
		if err := productConn.Close(); err != nil {
			log.Printf("failed to close product gRPC connection: %v", err)
		}
		if err := orderConn.Close(); err != nil {
			log.Printf("failed to close order gRPC connection: %v", err)
		}
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment gRPC connection: %v", err)
		}
	})

	cleanupFunc := func() {
		for _, fn := range cleanupFuncs {
			fn()
		}
		log.Println("gRPC clients cleaned up")
	}

	log.Println("gRPC clients initialized")

	return clients, cleanupFunc, nil
}
