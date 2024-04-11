package square

type Merchant struct {
	ID                        string
	TenantID                  string
	Name                      string
	SquareWebhookSignatureKey string
	SquareAPIToken            string
}

type MerchantRepository interface {
	PutSquareMerchant(merchant Merchant) error
	GetSquareMerchant(squareMerchantID string) (Merchant, error)
}
