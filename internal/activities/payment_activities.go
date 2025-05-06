package activities

import (
	"context"
	"fmt"
	"log"

	paymentpb "github.com/vogiaan1904/payment-svc/protogen/golang/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentActivities struct {
	PaymentSvcAddr string
}

func (a *PaymentActivities) ProcessPayment(ctx context.Context, request *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	conn, err := grpc.Dial(a.PaymentSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer conn.Close()

	client := paymentpb.NewPaymentServiceClient(conn)
	resp, err := client.ProcessPayment(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	log.Printf("Processed payment for order %s", request.OrderId)
	return resp, nil
}

func (a *PaymentActivities) GetPaymentStatus(ctx context.Context, paymentID string) (*paymentpb.GetPaymentStatusResponse, error) {
	conn, err := grpc.Dial(a.PaymentSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer conn.Close()

	client := paymentpb.NewPaymentServiceClient(conn)
	resp, err := client.GetPaymentStatus(ctx, &paymentpb.GetPaymentStatusRequest{
		PaymentId: paymentID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}

	return resp, nil
}

func (a *PaymentActivities) CancelPayment(ctx context.Context, paymentID string, reason string) error {
	conn, err := grpc.Dial(a.PaymentSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer conn.Close()

	client := paymentpb.NewPaymentServiceClient(conn)
	_, err = client.CancelPayment(ctx, &paymentpb.CancelPaymentRequest{
		PaymentId: paymentID,
		Reason:    reason,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	log.Printf("Cancelled payment %s", paymentID)
	return nil
}