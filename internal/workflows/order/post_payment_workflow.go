package order

import (
	"fmt"
	"time"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func ProcessPostPaymentOrder(ctx workflow.Context, params PostPaymentOrderWorkflowParams) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ProcessPostPaymentOrder started", "OrderCode", params.OrderCode)

	oData, err := validateOrder(ctx, params.OrderCode)
	if err != nil {
		return err
	}
	if oData.Status != orderpb.OrderStatus_PAYMENT_PENDING {
		logger.Info("Order already completed or cancelled", "OrderCode", params.OrderCode)
		return fmt.Errorf("order already completed or cancelled")
	}

	postPaymentCtx := workflow.WithActivityOptions(ctx, getPostPaymentActivityOptions())

	if err := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_PAYMENT_SUCCESS); err != nil {
		logger.Error("Failed to update order status to PAYMENT_SUCCESS, retrying...",
			"OrderCode", params.OrderCode, "Error", err)
		_ = workflow.Sleep(ctx, time.Second*5)
		if retryErr := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_PAYMENT_SUCCESS); retryErr != nil {
			return fmt.Errorf("failed to update order status after retry: %w", retryErr)
		}
	}

	if err := updateStock(postPaymentCtx, oData.Items); err != nil {
		logger.Error("Failed to update stock, creating manual intervention task",
			"OrderCode", params.OrderCode, "Error", err)
		// TODO: Implement manual intervention task creation
		// Continue workflow as stock update will be handled manually
	}

	// Deliver Order - Async step
	// TODO: Implement delivery workflow

	if err := updateOrderStatus(postPaymentCtx, params.OrderCode, orderpb.OrderStatus_COMPLETED); err != nil {
		logger.Error("Failed to mark order as completed, manual intervention required",
			"OrderCode", params.OrderCode, "Error", err)
		return fmt.Errorf("failed to complete order: %w", err)
	}

	logger.Info("Post-payment workflow completed successfully", "OrderCode", params.OrderCode)
	return nil
}
