package activities

import (
	"context"
	"fmt"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
)

type OrderActivities struct {
	Client orderpb.OrderServiceClient
}

func (a *OrderActivities) GetOrder(ctx context.Context, orderCode string) (*orderpb.OrderData, error) {
	resp, err := a.Client.FindOne(ctx, &orderpb.FindOneRequest{Request: &orderpb.FindOneRequest_Code{Code: orderCode}})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return resp.Order, nil
}

func (a *OrderActivities) UpdateOrderStatus(ctx context.Context, orderCode string, status orderpb.OrderStatus) error {
	_, err := a.Client.UpdateStatus(ctx, &orderpb.UpdateStatusRequest{
		Request: &orderpb.UpdateStatusRequest_Code{Code: orderCode},
		Status:  status,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
