package workflows

import (
	"time"

	"github.com/vogiaan1904/order-orchestrator/internal/activities"
	orderpb "github.com/vogiaan1904/order-svc/protogen/golang/order"
	paymentpb "github.com/vogiaan1904/payment-svc/protogen/golang/payment"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// OrderWorkflowParams contains the parameters for the order workflow
type OrderWorkflowParams struct {
	OrderID      string
	UserID       string
	Amount       float64
	PaymentMethod paymentpb.PaymentMethod
	Description  string
	Metadata     map[string]string
}

// OrderProcessingWorkflow coordinates the order and payment process
func OrderProcessingWorkflow(ctx workflow.Context, params OrderWorkflowParams) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("OrderProcessingWorkflow started", "OrderID", params.OrderID)

	// Activity options with automatic retries
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    30 * time.Second,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// 1. Get order details to validate
	var order *orderpb.Order
	err := workflow.ExecuteActivity(ctx, (*activities.OrderActivities).GetOrder, params.OrderID).Get(ctx, &order)
	if err != nil {
		logger.Error("Failed to get order", "Error", err)
		return "", err
	}

	// 2. Update order status to PAYMENT_PENDING
	err = workflow.ExecuteActivity(ctx, (*activities.OrderActivities).UpdateOrderStatus, 
		params.OrderID, orderpb.OrderStatus_ORDER_STATUS_PAYMENT_PENDING).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update order status to payment pending", "Error", err)
		return "", err
	}

	// 3. Process payment
	paymentRequest := &paymentpb.ProcessPaymentRequest{
		OrderId:     params.OrderID,
		UserId:      params.UserID,
		Amount:      params.Amount,
		Method:      params.PaymentMethod,
		Description: params.Description,
		Metadata:    params.Metadata,
	}

	var paymentResponse *paymentpb.ProcessPaymentResponse
	err = workflow.ExecuteActivity(ctx, (*activities.PaymentActivities).ProcessPayment, paymentRequest).Get(ctx, &paymentResponse)
	if err != nil {
		logger.Error("Payment processing failed", "Error", err)
		
		// Update order status to PAYMENT_FAILED
		updateErr := workflow.ExecuteActivity(ctx, (*activities.OrderActivities).UpdateOrderStatus,
			params.OrderID, orderpb.OrderStatus_ORDER_STATUS_PAYMENT_FAILED).Get(ctx, nil)
		if updateErr != nil {
			logger.Error("Failed to update order status to payment failed", "Error", updateErr)
		}
		
		return "", err
	}

	// 4. Update order status to PAID
	err = workflow.ExecuteActivity(ctx, (*activities.OrderActivities).UpdateOrderStatus,
		params.OrderID, orderpb.OrderStatus_ORDER_STATUS_PAID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to update order status to paid", "Error", err)
		return "", err
	}

	// We need the payment URL for redirecting the user
	return paymentResponse.PaymentUrl, nil
}