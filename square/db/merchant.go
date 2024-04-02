package db

type SquareMerchant struct {
	SquareMerchantID          string
	SquareWebhookSignatureKey string
	SquareAPIKey              string
}

type SquareMerchantsRepository interface {
	CreateMerchant(merchant SquareMerchant) error
	FindById(squareMerchantID string) (SquareMerchant, error)
}
