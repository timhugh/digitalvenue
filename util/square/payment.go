package square

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type PaymentRepository interface {
	PutSquarePayment(payment *Payment) error
	GetSquarePayment(squareMerchantID string, squarePaymentID string) (*Payment, error)
}
