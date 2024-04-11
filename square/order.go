package square

import (
	"fmt"
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

func MapOrder(squareOrder *Order, squarePaymentID string, squareMerchantID string, tenantID string, customerID string) (*core.Order, error) {
	order := core.Order{
		ID:         squareOrder.SquareOrderID, // We use the SquareOrderID as the ID for orders that come from Square to make duplicate detection easier
		TenantID:   tenantID,
		CustomerID: customerID,
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
			return nil, err
		}

		for i := 0; i < quantity; i++ {
			order.Items = append(order.Items, core.OrderItem{
				Name: item.Name,
				ID:   fmt.Sprintf("%s-%d", item.ItemID, i+1),
				Meta: map[string]string{
					ItemIDKey: item.ItemID,
				},
			})
		}
	}

	return &order, nil
}
