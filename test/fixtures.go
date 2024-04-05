package test

import (
	"github.com/timhugh/digitalvenue/core"
)

const (
	CustomerID        = "customer_id"
	CustomerFirstName = "Tim"
	CustomerLastName  = "Heuett"
	CustomerEmail     = "info@timheuett.com"
	CustomerPhone     = "+12062062062"

	OrderID   = "order_id"
	ItemName1 = "Item 1"
	ItemName2 = "Item 2"
)

func NewOrder() core.Order {
	return core.Order{
		Items: []core.OrderItem{
			{
				Name: ItemName1,
			},
			{
				Name: ItemName2,
			},
			{
				Name: ItemName2,
			},
		},
	}
}

func NewCustomer() core.Customer {
	return core.Customer{
		FirstName: CustomerFirstName,
		LastName:  CustomerLastName,
		Email:     CustomerEmail,
		Phone:     CustomerPhone,
	}
}
