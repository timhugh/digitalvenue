package square

type Merchant struct {
	SquareMerchantID          string
	SquareWebhookSignatureKey string
	SquareAPIToken            string
}

type MerchantsRepository interface {
	Put(merchant Merchant) error
	Get(squareMerchantID string) (Merchant, error)
}
