package square

import (
	"fmt"
	"github.com/timhugh/digitalvenue/core"
	"strconv"
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

type orderContainer struct {
	Order Order `json:"order"`
}

func (client client) GetOrder(squareOrderID string, apiToken string) (Order, error) {
	path := client.baseUrl + fmt.Sprintf(getOrderRouteFormat, squareOrderID)

	var orderContainer orderContainer
	err := client.fetchJson(path, apiToken, &orderContainer)
	if err != nil {
		return Order{}, err
	}

	return orderContainer.Order, nil
}

func MapOrder(squareOrder Order) (core.Order, error) {
	order := core.Order{
		Meta: core.OrderMeta{
			SquareOrderID:    squareOrder.SquareOrderID,
			SquareCustomerID: squareOrder.SquareCustomerID,
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
				Meta: core.OrderItemMeta{
					SquareItemID: item.ItemID,
				},
			})
		}
	}

	return order, nil
}
