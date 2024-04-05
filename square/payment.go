package square

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type PaymentsRepository interface {
	Put(payment Payment) error
	Get(squarePaymentID string) (Payment, error)
}

type PaymentCreatedEvent struct {
	SquarePaymentId string
}

type PaymentCreatedQueue interface {
	Publish(squarePaymentId string) error
}
