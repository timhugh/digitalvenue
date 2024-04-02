package db

type SquarePayment struct {
	SquarePaymentID  string
	SquareMerchantID string
	SquareOrderID    string
}

type SquarePaymentsRepository interface {
	Create(payment SquarePayment) error
}
