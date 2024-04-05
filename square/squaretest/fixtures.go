package squaretest

import (
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/test"
)

const (
	SquareMerchantID = "squareMerchantID"
	SquarePaymentID  = "squarePaymentID"
	SquareCustomerID = "squareCustomerID"
	SquareLocationID = "squareLocationID"
	SquareOrderID    = "squareOrderID"

	SquareItemUID1 = "squareItemUID1"
	SquareItemUID2 = "squareItemUID2"
	SquareAPIToken = "squareAPIToken"
)

func NewSquarePayment() square.Payment {
	return square.Payment{
		SquareMerchantID: SquareMerchantID,
		SquareOrderID:    SquareOrderID,
		SquarePaymentID:  SquarePaymentID,
	}
}

func NewSquareMerchant() square.Merchant {
	return square.Merchant{
		SquareMerchantID: SquareMerchantID,
		SquareAPIToken:   SquareAPIToken,
	}
}

func NewSquareOrder() square.Order {
	return square.Order{
		SquareOrderID:    SquareOrderID,
		SquareCustomerID: SquareCustomerID,
		SquareLocationID: SquareLocationID,
		OrderItems: []square.OrderItem{
			{
				ItemID:   SquareItemUID1,
				Name:     test.ItemName1,
				Quantity: "1",
			},
			{
				ItemID:   SquareItemUID2,
				Name:     test.ItemName2,
				Quantity: "2",
			},
		},
	}
}

func NewSquareCustomer() square.Customer {
	return square.Customer{
		SquareCustomerID: SquareCustomerID,
		FirstName:        test.CustomerFirstName,
		LastName:         test.CustomerLastName,
		Email:            test.CustomerEmail,
		Phone:            test.CustomerPhone,
	}
}

func NewOrder() core.Order {
	order := test.NewOrder()
	order.Meta.SquareOrderID = SquareOrderID
	order.Meta.SquareCustomerID = SquareCustomerID

	order.Items[0].Meta.SquareItemID = SquareItemUID1
	order.Items[1].Meta.SquareItemID = SquareItemUID2
	order.Items[2].Meta.SquareItemID = SquareItemUID2

	return order
}

func NewOrderGathered() core.Order {
	order := NewOrder()

	order.CustomerID = test.CustomerID

	order.Meta.SquareMerchantID = SquareMerchantID
	order.Meta.SquarePaymentID = SquarePaymentID
	order.Meta.SquareCustomerID = SquareCustomerID

	return order
}

func NewCustomer() core.Customer {
	customer := test.NewCustomer()
	customer.Meta.SquareCustomerID = SquareCustomerID
	return customer
}
