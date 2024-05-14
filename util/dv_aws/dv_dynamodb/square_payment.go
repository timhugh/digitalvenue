package dv_dynamodb

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/square"
)

type squarePaymentDTO struct {
	PK            string
	SK            string
	TenantID      string
	SquareOrderID string
}

func (repo *Repository) GetSquarePayment(squareMerchantID string, squarePaymentID string) (*square.Payment, error) {
	key := map[string]string{
		"PK": PrefixID("SquareMerchant", squareMerchantID),
		"SK": PrefixID("SquarePayment", squarePaymentID),
	}

	item := squarePaymentDTO{}
	err := repo.get("SquarePayment", key, &item)
	if err != nil {
		if errors.Is(err, ItemNotFoundException{}) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get SquarePayment")
	}

	return &square.Payment{
		TenantID:         item.TenantID,
		SquarePaymentID:  squarePaymentID,
		SquareMerchantID: squareMerchantID,
		SquareOrderID:    item.SquareOrderID,
	}, nil
}

func (repo *Repository) PutSquarePayment(payment *square.Payment) error {
	inputPayment := &squarePaymentDTO{
		PK:            PrefixID("SquareMerchant", payment.SquareMerchantID),
		SK:            PrefixID("SquarePayment", payment.SquarePaymentID),
		TenantID:      PrefixID("Tenant", payment.TenantID),
		SquareOrderID: payment.SquareOrderID,
	}

	err := repo.put("SquarePayment", inputPayment)
	if err != nil {
		return errors.Wrap(err, "failed to put SquarePayment")
	}

	return nil
}
