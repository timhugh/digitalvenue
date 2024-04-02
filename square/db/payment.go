package db

type Payment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type PaymentsRepository interface {
	CreatePayment(payment Payment) error
}
