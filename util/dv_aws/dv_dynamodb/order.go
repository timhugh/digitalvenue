package dv_dynamodb

import (
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"maps"
)

type orderDTO struct {
	PK         string
	SK         string
	CustomerID string
	Meta       map[string]string
	OrderItems []orderItemDTO
}

type orderItemDTO struct {
	ItemID string
	Name   string
	Meta   map[string]string
}

func (repo *Repository) GetOrder(tenantID string, orderID string) (*core.Order, error) {
	key := map[string]string{
		"PK": PrefixID("Tenant", tenantID),
		"SK": PrefixID("Order", orderID),
	}

	item := orderDTO{}
	err := repo.get("Order", key, &item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Order")
	}

	orderItems := make([]core.OrderItem, len(item.OrderItems))
	for i, item := range item.OrderItems {
		orderItems[i] = core.OrderItem{
			ID:   item.ItemID,
			Name: item.Name,
			Meta: maps.Clone(item.Meta),
		}
	}

	return &core.Order{
		ID:         orderID,
		TenantID:   tenantID,
		CustomerID: item.CustomerID,
		Items:      orderItems,
		Meta:       maps.Clone(item.Meta),
	}, nil
}

func (repo *Repository) PutOrder(order *core.Order) error {
	items := make([]orderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = orderItemDTO{
			ItemID: item.ID,
			Name:   item.Name,
			Meta:   maps.Clone(item.Meta),
		}
	}

	inputOrder := &orderDTO{
		PK:         PrefixID("Tenant", order.TenantID),
		SK:         PrefixID("Order", order.ID),
		CustomerID: order.CustomerID,
		Meta:       maps.Clone(order.Meta),
		OrderItems: items,
	}

	err := repo.put("Order", inputOrder)
	if err != nil {
		return errors.Wrap(err, "failed to put order")
	}

	return nil
}
