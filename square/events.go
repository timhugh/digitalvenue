package square

type PaymentCreated struct {
	SquareMerchantID string
	SquarePaymentID  string
}

type CustomerCreated struct {
	SquareMerchantID string
	SquareCustomerID string
}
