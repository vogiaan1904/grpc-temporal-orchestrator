package order

import (
	"fmt"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func ProcessPrePaymentOrder(ctx workflow.Context, params OrderWorkflowParams) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ProcessPrePaymentOrder started", "OrderID", params.OrderID)

	// 1. Validate Order
	oData, err := validateOrder(ctx, params.OrderID)
	if err != nil {
		return "", err
	}

	// 2. Reserve Inventory
	if err := reserveInventory(ctx, oData.Items); err != nil {
		return "", err
	}

	// 3. Update to PAYMENT_PENDING
	if err := updateOrderStatus(ctx, params.OrderID, orderpb.OrderStatus_PAYMENT_PENDING); err != nil {
		// Compensate: Release inventory
		if compensationErr := releaseInventory(ctx, oData.Items); compensationErr != nil {
			logger.Error("Compensation failed", "Error", compensationErr)
		}

		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	// 4. Process Payment
	paymentResponse, err := processPayment(ctx, params)
	if err != nil {
		// Compensate: Release inventory
		if compensationErr := releaseInventory(ctx, oData.Items); compensationErr != nil {
			logger.Error("Compensation failed", "Error", compensationErr)
		}
		return "", fmt.Errorf("payment failed: %w", err)
	}

	workflow.Go(ctx, func(ctx workflow.Context) {
		if err := workflow.Sleep(ctx, PaymentTimeout); err != nil {
			return // Workflow context canceled
		}

		o, err := validateOrder(ctx, params.OrderID)
		if err == nil && o.Status == orderpb.OrderStatus_PAYMENT_PENDING {
			_ = updateOrderStatus(ctx, params.OrderID, orderpb.OrderStatus_PAYMENT_FAILED)
			_ = releaseInventory(ctx, oData.Items)
		}
	})

	return paymentResponse.PaymentUrl, nil
}
