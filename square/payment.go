package square

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
	TenantID         string
}

type PaymentRepository interface {
	PutSquarePayment(payment Payment) error
	GetSquarePayment(squarePaymentID string) (Payment, error)
}
