package order

import (
	"fmt"

	"github.com/vogiaan1904/order-orchestrator/internal/activities"
	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	paymentpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/payment"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func validateOrder(ctx workflow.Context, orderCode string) (*orderpb.OrderData, error) {
	var orderData *orderpb.OrderData
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getOrderActivityOptions()),
		(*activities.OrderActivities).GetOrder,
		orderCode,
	).Get(ctx, &orderData)
	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(
			"Failed to get order",
			"ORDER_NOT_FOUND",
			err,
		)
	}
	return orderData, nil
}

func reserveInventory(ctx workflow.Context, items []*orderpb.OrderItem) error {
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getProductActivityOptions()),
		(*activities.ProductActivities).ReserveInventory,
		items,
	).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to reserve inventory: %w", err)
	}
	return nil
}

func updateOrderStatus(ctx workflow.Context, orderID string, status orderpb.OrderStatus) error {
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getOrderActivityOptions()),
		(*activities.OrderActivities).UpdateOrderStatus,
		orderID,
		status,
	).Get(ctx, nil)
	return err
}

func processPayment(ctx workflow.Context, params PrePaymentOrderWorkflowParams) (*paymentpb.ProcessPaymentResponse, error) {
	paymentRequest := &paymentpb.ProcessPaymentRequest{
		OrderCode:   params.OrderCode,
		UserId:      params.UserID,
		Amount:      params.Amount,
		Method:      params.PaymentMethod,
		Description: params.Description,
		Metadata:    params.Metadata,
	}

	var paymentResponse *paymentpb.ProcessPaymentResponse
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getPaymentActivityOptions()),
		(*activities.PaymentActivities).ProcessPayment,
		paymentRequest,
	).Get(ctx, &paymentResponse)
	return paymentResponse, err
}

func updateStock(ctx workflow.Context, items []*orderpb.OrderItem) error {
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, getProductActivityOptions()),
		(*activities.ProductActivities).UpdateStock,
		items,
	).Get(ctx, nil)
	return err
}
