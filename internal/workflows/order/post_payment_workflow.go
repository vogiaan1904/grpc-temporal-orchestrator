package order

import (
	"fmt"
	"time"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func ProcessPostPaymentOrder(ctx workflow.Context, params PostPaymentOrderWorkflowParams) error {
	oData, err := validateOrder(ctx, params.OrderCode)
	if err != nil {
		return err
	}
	if oData.Status != orderpb.OrderStatus_PAYMENT_PENDING {
		return fmt.Errorf("order already completed or cancelled")
	}

	postPaymentCtx := workflow.WithActivityOptions(ctx, getPostPaymentActivityOptions())

	if err := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_PAYMENT_SUCCESS); err != nil {
		_ = workflow.Sleep(ctx, time.Second*5)
		if retryErr := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_PAYMENT_SUCCESS); retryErr != nil {
			return fmt.Errorf("failed to update order status after retry: %w", retryErr)
		}
	}

	if err := updateStock(postPaymentCtx, oData.Items); err != nil {
		// TODO: Implement manual intervention task creation
		// Continue workflow as stock update will be handled manually
	}

	// Deliver Order - Async step
	// TODO: Implement delivery workflow

	if err := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_COMPLETED); err != nil {
		return fmt.Errorf("failed to complete order: %w", err)
	}

	return nil
}
