package core

type Order struct {
	ID         string
	TenantID   string
	CustomerID string
	Items      []OrderItem
	Meta       map[string]string
}

type OrderItem struct {
	ID   string
	Name string
	Meta map[string]string
}

type OrderRepository interface {
	PutOrder(order *Order) error
}
