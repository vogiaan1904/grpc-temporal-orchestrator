package order

import (
	"fmt"

	orderpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/order"
	"go.temporal.io/sdk/workflow"
)

func ProcessPrePaymentOrder(ctx workflow.Context, params PrePaymentOrderWorkflowParams) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("ProcessPrePaymentOrder started", "OrderCode", params.OrderCode)

	oData, err := validateOrder(ctx, params.OrderCode)
	if err != nil {
		return "", fmt.Errorf("failed to validate order: %w", err)
	}

	if err := reserveInventory(ctx, oData.Items); err != nil {
		return "", fmt.Errorf("failed to reserve inventory: %w", err)
	}

	if err := updateOrderStatus(ctx, params.OrderCode, orderpb.OrderStatus_PAYMENT_PENDING); err != nil {
		if compensationErr := releaseInventory(ctx, oData.Items); compensationErr != nil {
			logger.Error("Compensation failed", "Error", compensationErr)
		}

		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	paymentResponse, err := processPayment(ctx, params)
	if err != nil {
		if compensationErr := releaseInventory(ctx, oData.Items); compensationErr != nil {
			logger.Error("Compensation failed", "Error", compensationErr)
		}

		return "", fmt.Errorf("payment failed: %w", err)
	}

	workflow.Go(ctx, func(ctx workflow.Context) {
		if err := workflow.Sleep(ctx, PaymentTimeout); err != nil {
			return
		}

		o, err := validateOrder(ctx, params.OrderCode)
		if err == nil && o.Status == orderpb.OrderStatus_PAYMENT_PENDING {
			_ = updateOrderStatus(ctx, params.OrderCode, orderpb.OrderStatus_PAYMENT_FAILED)
			_ = releaseInventory(ctx, oData.Items)
		}
	})

	return paymentResponse.PaymentUrl, nil
}
