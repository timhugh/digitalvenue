package core

type Order struct {
	OrderID    string
	TenantID   string
	CustomerID string
	Items      []OrderItem
	Meta       map[string]string
}

type OrderItem struct {
	ItemID string
	Name   string
	Meta   map[string]string
}

type OrderRepository interface {
	PutOrder(order Order) (string, error)
}
