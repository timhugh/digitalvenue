package square

import (
	"github.com/timhugh/digitalvenue/core"
	"strconv"
)

type OrderMapper interface {
	MapOrder(squareOrder Order) (core.Order, error)
}

type orderMapper struct{}

func NewOrderMapper() OrderMapper {
	return &orderMapper{}
}

func (o *orderMapper) MapOrder(squareOrder Order) (core.Order, error) {
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
