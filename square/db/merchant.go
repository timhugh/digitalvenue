package db

type Merchant struct {
	SquareMerchantId          string
	SquareWebhookSignatureKey string
	SquareAPIKey              string
}

type MerchantsRepository interface {
	CreateMerchant(merchant Merchant) error
	FindMerchantBySquareMerchantId(squareMerchantId string) (Merchant, error)
}
