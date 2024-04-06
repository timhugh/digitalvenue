package square

import (
	"github.com/timhugh/digitalvenue/core"
	"strconv"
)

const (
	OrderIDKey    = "SquareOrderID"
	ItemIDKey     = "SquareItemID"
	MerchantIDKey = "SquareMerchantID"
	PaymentIDKey  = "SquarePaymentID"
)

type Order struct {
	SquareOrderID    string      `json:"id"`
	SquareCustomerID string      `json:"customer_id"`
	SquareLocationID string      `json:"location_id"`
	OrderItems       []OrderItem `json:"line_items"`
}

type OrderItem struct {
	ItemID   string `json:"uid"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

func MapOrder(squareOrder Order, squarePaymentID string, squareMerchantID string) (core.Order, error) {
	order := core.Order{
		Meta: map[string]string{
			MerchantIDKey: squareMerchantID,
			PaymentIDKey:  squarePaymentID,
			OrderIDKey:    squareOrder.SquareOrderID,
			CustomerIDKey: squareOrder.SquareCustomerID,
		},
	}

	for _, item := range squareOrder.OrderItems {
		quantity, err := strconv.Atoi(item.Quantity)
		if err != nil {
			return order, err
		}
		for i := 0; i < quantity; i++ {
			order.Items = append(order.Items, core.OrderItem{
				Name: item.Name,
				Meta: map[string]string{
					ItemIDKey: item.ItemID,
				},
			})
		}
	}

	return order, nil
}
