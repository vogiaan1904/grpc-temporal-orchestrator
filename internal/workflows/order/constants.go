package order

import "time"

const (
	PaymentTimeout             = 6 * time.Minute //( 5 minutes + 1 minute for callback)
	SignalNamePaymentCompleted = "payment-completed"
)
