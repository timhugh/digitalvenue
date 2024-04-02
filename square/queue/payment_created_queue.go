package queue

type SquarePaymentCreatedEvent struct {
	SquarePaymentId string
}

type SquarePaymentCreatedQueue interface {
	Publish(squarePaymentId string) error
}
