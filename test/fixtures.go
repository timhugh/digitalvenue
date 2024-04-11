package test

import (
	"github.com/timhugh/digitalvenue/core"
)

const (
	TenantID   = "tenantID"
	TenantName = "Test Tenant"

	CustomerID        = "customerID"
	CustomerName      = "Tim Heuett"
	CustomerFirstName = "Tim"
	CustomerLastName  = "Heuett"
	CustomerEmail     = "info@timheuett.com"
	CustomerPhone     = "+12062062062"

	OrderID   = "orderID"
	ItemID1   = "itemID1"
	ItemName1 = "Item 1"
	ItemID2   = "itemID2"
	ItemName2 = "Item 2"
	ItemID3   = "itemID3"
)

func NewOrder() *core.Order {
	return &core.Order{
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

func NewCustomer() *core.Customer {
	return &core.Customer{
		TenantID: TenantID,
		ID:       CustomerID,
		Name:     CustomerName,
		Email:    CustomerEmail,
		Phone:    CustomerPhone,
	}
}
