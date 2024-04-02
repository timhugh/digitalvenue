package queue

type PaymentCreatedEvent struct {
	SquarePaymentId string
}

type PaymentCreatedQueue interface {
	Publish(squarePaymentId string) error
}
