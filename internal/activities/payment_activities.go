package activities

import (
	"context"
	"fmt"

	paymentpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/payment"
	"google.golang.org/grpc"
)

type PaymentActivities struct {
	Client paymentpb.PaymentServiceClient
}

func NewPaymentActivities(conn *grpc.ClientConn) *PaymentActivities {
	return &PaymentActivities{
		Client: paymentpb.NewPaymentServiceClient(conn),
	}
}

func (a *PaymentActivities) ProcessPayment(ctx context.Context, request *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	resp, err := a.Client.ProcessPayment(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	return resp, nil
}

func (a *PaymentActivities) CancelPayment(ctx context.Context, orderCode string, reason string) error {
	_, err := a.Client.CancelPayment(ctx, &paymentpb.CancelPaymentRequest{
		PaymentIdentifier: &paymentpb.CancelPaymentRequest_OrderCode{
			OrderCode: orderCode,
		},
		Reason: reason,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	return nil
}
