package db

type Merchant struct {
	SquareMerchantID          string
	SquareWebhookSignatureKey string
	SquareAPIKey              string
}

type MerchantsRepository interface {
	CreateMerchant(merchant Merchant) error
	FindMerchantBySquareMerchantID(squareMerchantID string) (Merchant, error)
}
