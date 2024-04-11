package test

import (
	"github.com/timhugh/digitalvenue/core"
)

const (
	TenantID   = "tenant_id"
	TenantName = "Test Tenant"

	CustomerID        = "customer_id"
	CustomerName      = "Tim Heuett"
	CustomerFirstName = "Tim"
	CustomerLastName  = "Heuett"
	CustomerEmail     = "info@timheuett.com"
	CustomerPhone     = "+12062062062"

	OrderID   = "order_id"
	ItemID1   = "item_id_1"
	ItemName1 = "Item 1"
	ItemID2   = "item_id_2"
	ItemName2 = "Item 2"
	ItemID3   = "item_id_3"
)

func NewOrder() core.Order {
	return core.Order{
		ID:         OrderID,
		TenantID:   TenantID,
		CustomerID: CustomerID,
		Items: []core.OrderItem{
			{
				ID:   ItemID1,
				Name: ItemName1,
			},
			{
				ID:   ItemID2,
				Name: ItemName2,
			},
			{
				ID:   ItemID3,
				Name: ItemName2,
			},
		},
	}
}

func NewCustomer() core.Customer {
	return core.Customer{
		TenantID: TenantID,
		ID:       CustomerID,
		Name:     CustomerName,
		Email:    CustomerEmail,
		Phone:    CustomerPhone,
	}
}
