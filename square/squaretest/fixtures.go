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

	SquareItemUID1            = "squareItemUID1"
	SquareItemUID2            = "squareItemUID2"
	SquareAPIToken            = "squareAPIToken"
	SquareWebhookSignatureKey = "squareWebhookSignatureKey"
)

func NewSquarePayment() square.Payment {
	return square.Payment{
		SquareMerchantID: SquareMerchantID,
		SquareOrderID:    SquareOrderID,
		SquarePaymentID:  SquarePaymentID,
		TenantID:         test.TenantID,
	}
}

func NewSquareMerchant() square.Merchant {
	return square.Merchant{
		ID:                        SquareMerchantID,
		TenantID:                  test.TenantID,
		Name:                      test.TenantName,
		SquareAPIToken:            SquareAPIToken,
		SquareWebhookSignatureKey: SquareWebhookSignatureKey,
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
	order.Meta = map[string]string{
		square.OrderIDKey:    SquareOrderID,
		square.CustomerIDKey: SquareCustomerID,
		square.PaymentIDKey:  SquarePaymentID,
		square.MerchantIDKey: SquareMerchantID,
	}

	order.Items[0].Meta = map[string]string{square.ItemIDKey: SquareItemUID1}
	order.Items[1].Meta = map[string]string{square.ItemIDKey: SquareItemUID2}
	order.Items[2].Meta = map[string]string{square.ItemIDKey: SquareItemUID2}

	return order
}

func NewCustomer() core.Customer {
	customer := test.NewCustomer()
	customer.Meta = map[string]string{
		"SquareCustomerID": SquareCustomerID,
	}
	return customer
}
