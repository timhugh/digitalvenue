package square

type Merchant struct {
	ID       string
	TenantID string
	Name     string

	SquareWebhookSignatureKey string
	SquareAPIToken            string

	TicketableCategories []string
}

type MerchantRepository interface {
	GetSquareMerchant(squareMerchantID string) (*Merchant, error)
}
