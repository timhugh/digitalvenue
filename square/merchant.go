package square

type Merchant struct {
	SquareMerchantID          string
	SquareWebhookSignatureKey string
	SquareAPIToken            string
}

type MerchantRepository interface {
	PutSquareMerchant(merchant Merchant) error
	GetSquareMerchant(squareMerchantID string) (Merchant, error)
}
