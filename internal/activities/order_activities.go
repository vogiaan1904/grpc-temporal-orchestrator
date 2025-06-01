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

func (a *OrderActivities) GetOrder(ctx context.Context, orderID string) (*orderpb.OrderData, error) {
	resp, err := a.Client.FindOne(ctx, &orderpb.FindOneRequest{Id: orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	log.Printf("Retrieved order %s: %v", orderID, resp.Order)
	return resp.Order, nil
}

func (a *OrderActivities) UpdateOrderStatus(ctx context.Context, orderID string, status orderpb.OrderStatus) error {
	_, err := a.Client.UpdateStatus(ctx, &orderpb.UpdateStatusRequest{
		Id:     orderID,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Updated order %s status to %s", orderID, status.String())
	return nil
}
