package square

type Merchant struct {
	SquareMerchantID          string
	SquareWebhookSignatureKey string
	SquareAPIToken            string
}

type MerchantsRepository interface {
	Create(merchant Merchant) error
	FindByID(squareMerchantID string) (Merchant, error)
}