package dv_dynamodb

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/square"
)

type squareMerchant struct {
	PK                        string
	SK                        string
	Type                      string
	TenantID                  string
	Name                      string
	SquareAPIToken            string
	SquareWebhookSignatureKey string
}

func (repo *Repository) GetSquareMerchant(squareMerchantID string) (*square.Merchant, error) {
	merchantKey := PrefixID("SquareMerchant", squareMerchantID)
	key := map[string]string{
		"PK": merchantKey,
		"SK": merchantKey,
	}

	item := squareMerchant{}
	err := repo.get("SquareMerchant", key, &item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get SquareMerchant")
	}

	tenantID, err := UnprefixID(item.TenantID)
	if err != nil {
		return nil, errors.Wrap(err, "SquareMerchant has invalid TenantID")
	}

	return &square.Merchant{
		TenantID:                  tenantID,
		Name:                      item.Name,
		ID:                        squareMerchantID,
		SquareWebhookSignatureKey: item.SquareWebhookSignatureKey,
		SquareAPIToken:            item.SquareAPIToken,
	}, nil
}
