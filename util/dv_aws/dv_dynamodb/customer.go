package dv_dynamodb

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"maps"
)

type customerDTO struct {
	PK         string
	SK         string
	CustomerID string
	Name       string
	Email      string
	Phone      string
	Meta       map[string]string
}

func (repo *Repository) GetCustomer(tenantID string, customerID string) (*core.Customer, error) {
	key := map[string]string{
		"PK": PrefixID("Tenant", tenantID),
		"SK": PrefixID("Customer", customerID),
	}

	item := customerDTO{}
	err := repo.get("Customer", key, &item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Customer")
	}

	return &core.Customer{
		TenantID: tenantID,
		ID:       customerID,
		Name:     item.Name,
		Email:    item.Email,
		Phone:    item.Phone,
		Meta:     maps.Clone(item.Meta),
	}, nil
}

func (repo *Repository) PutCustomer(customer *core.Customer) error {
	inputCustomer := &customerDTO{
		PK:         PrefixID("Tenant", customer.TenantID),
		SK:         PrefixID("Customer", customer.ID),
		CustomerID: customer.ID,
		Name:       customer.Name,
		Email:      customer.Email,
		Phone:      customer.Phone,
		Meta:       maps.Clone(customer.Meta),
	}

	err := repo.put("Customer", inputCustomer)
	if err != nil {
		return errors.Wrap(err, "failed to put Customer")
	}

	return nil
}
