package square

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type PaymentRepository interface {
	PutSquarePayment(payment Payment) error
	GetSquarePayment(squarePaymentID string) (Payment, error)
}

type PaymentCreatedEvent struct {
	SquarePaymentId string
}

type PaymentCreatedQueue interface {
	PublishSquarePaymentCreated(squarePaymentId string) error
}
