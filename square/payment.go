package square

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type PaymentsRepository interface {
	Create(payment Payment) error
	FindByID(squarePaymentID string) (Payment, error)
}
