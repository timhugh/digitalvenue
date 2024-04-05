package core

type Order struct {
	OrderID    string
	TenantID   string
	CustomerID string
	Items      []OrderItem
	Meta       OrderMeta
}

type OrderMeta struct {
	SquareOrderID    string
	SquarePaymentID  string
	SquareMerchantID string
	SquareCustomerID string
}

type OrderItem struct {
	ItemID string
	Name   string
	Meta   OrderItemMeta
}

type OrderItemMeta struct {
	SquareItemID string
}

type OrderRepository interface {
	Put(order Order) (string, error)
}
