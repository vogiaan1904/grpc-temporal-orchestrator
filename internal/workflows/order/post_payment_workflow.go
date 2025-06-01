package order

import (
	"fmt"
	"time"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func ProcessPostPaymentOrder(ctx workflow.Context, params OrderWorkflowParams) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("ProcessPostPaymentOrder started", "OrderID", params.OrderID)

	// 1. Validate Order with standard options
	oData, err := validateOrder(ctx, params.OrderID)
	if err != nil {
		return err
	}
	if oData.Status != orderpb.OrderStatus_PAYMENT_PENDING {
		logger.Info("Order already completed or cancelled", "OrderID", params.OrderID)
		return fmt.Errorf("order already completed or cancelled")
	}

	postPaymentCtx := workflow.WithActivityOptions(ctx, getPostPaymentActivityOptions())

	// 2. Update to PAYMENT_SUCCESS - Critical step
	if err := updateOrderStatus(postPaymentCtx, params.OrderID, orderpb.OrderStatus_PAYMENT_SUCCESS); err != nil {
		logger.Error("Failed to update order status to PAYMENT_SUCCESS, retrying...",
			"OrderID", params.OrderID, "Error", err)
		_ = workflow.Sleep(ctx, time.Second*5)
		if retryErr := updateOrderStatus(postPaymentCtx, params.OrderID, orderpb.OrderStatus_PAYMENT_SUCCESS); retryErr != nil {
			return fmt.Errorf("failed to update order status after retry: %w", retryErr)
		}
	}

	// 3. Update Stock - Critical step
	if err := updateStock(postPaymentCtx, oData.Items); err != nil {
		logger.Error("Failed to update stock, creating manual intervention task",
			"OrderID", params.OrderID, "Error", err)
		// TODO: Implement manual intervention task creation
		// Continue workflow as stock update will be handled manually
	}

	// 4. Deliver Order - Async step
	// Create a separate child workflow for delivery tracking
	// childID := fmt.Sprintf("delivery-%s", params.OrderID)
	// cwo := workflow.ChildWorkflowOptions{
	// 	WorkflowID: childID,
	// 	// Add appropriate retry policy for delivery workflow
	// }
	// ctx = workflow.WithChildOptions(ctx, cwo)
	// TODO: Implement delivery workflow

	// 5. Update Order Status to COMPLETED - Final step
	if err := updateOrderStatus(postPaymentCtx, params.OrderID, orderpb.OrderStatus_COMPLETED); err != nil {
		logger.Error("Failed to mark order as completed, manual intervention required",
			"OrderID", params.OrderID, "Error", err)
		// Create a manual intervention task
		// TODO: Implement manual intervention task creation
		return fmt.Errorf("failed to complete order: %w", err)
	}

	logger.Info("Post-payment workflow completed successfully", "OrderID", params.OrderID)
	return nil
}
