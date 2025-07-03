package order

type PrePaymentOrderWorkflowParams struct {
	OrderCode       string
	UserID          string
	TotalAmount     float64
	Provider        string
	ProviderDetails string
	Metadata        map[string]string
}

type PostPaymentOrderWorkflowParams struct {
	OrderCode string
	Metadata  map[string]string
}
