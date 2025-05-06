package activities

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderActivities struct {
	OrderSvcAddr string
}

func (a *OrderActivities) GetOrder(ctx context.Context, orderID string) (*orderpb.Order, error) {
	conn, err := grpc.Dial(a.OrderSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)
	resp, err := client.GetOrder(ctx, &orderpb.GetOrderRequest{Id: orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	log.Printf("Retrieved order %s: %v", orderID, resp.Order)
	return resp.Order, nil
}

func (a *OrderActivities) UpdateOrderStatus(ctx context.Context, orderID string, status orderpb.OrderStatus) error {
	conn, err := grpc.Dial(a.OrderSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to order service: %w", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)
	_, err = client.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  status,
	})
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	log.Printf("Updated order %s status to %s", orderID, status.String())
	return nil
}
