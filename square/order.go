package square

import (
	"fmt"
)

type Order struct {
	SquareOrderID    string
	SquareCustomerID string
}

type OrdersRepository interface {
	Create(order Order) error
}

type orderContainer struct {
	Order struct {
		Id         string `json:"id"`
		CustomerId string `json:"customer_id"`
	} `json:"order"`
}

func (client client) GetOrder(squareOrderID string, apiToken string) (Order, error) {
	path := client.baseUrl + fmt.Sprintf(getOrderRouteFormat, squareOrderID)

	var orderContainer orderContainer
	err := client.fetchJson(path, apiToken, &orderContainer)
	if err != nil {
		return Order{}, err
	}

	return Order{
		SquareOrderID:    orderContainer.Order.Id,
		SquareCustomerID: orderContainer.Order.CustomerId,
	}, nil
}
