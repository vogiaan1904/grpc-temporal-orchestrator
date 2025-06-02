package order

import paymentpb "github.com/vogiaan1904/order-orchestrator/protogen/golang/payment"

// OrderWorkflowParams contains the parameters for the order workflow
type PrePaymentOrderWorkflowParams struct {
	OrderCode     string
	UserID        string
	Amount        float64
	PaymentMethod paymentpb.PaymentMethod
	Description   string
	Metadata      map[string]string
}

type PostPaymentOrderWorkflowParams struct {
	OrderCode string
	Metadata  map[string]string
}
