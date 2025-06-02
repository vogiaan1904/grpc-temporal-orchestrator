package activities

import (
	"context"
	"fmt"
	"log"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
)

type OrderActivities struct {
	Client orderpb.OrderServiceClient
}

func (a *OrderActivities) GetOrder(ctx context.Context, orderCode string) (*orderpb.OrderData, error) {
	resp, err := a.Client.FindOne(ctx, &orderpb.FindOneRequest{Code: orderCode})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	log.Printf("Retrieved order %s: %v", orderCode, resp.Order)
	return resp.Order, nil
}

func (a *OrderActivities) UpdateOrderStatus(ctx context.Context, orderCode string, status orderpb.OrderStatus) error {
	_, err := a.Client.UpdateStatus(ctx, &orderpb.UpdateStatusRequest{
		Code:   orderCode,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Updated order %s status to %s", orderCode, status.String())
	return nil
}
