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
	GetOrder(tenantID string, orderID string) (*Order, error)
	PutOrder(order *Order) error
}

type OrderCreatedQueue interface {
	PublishOrderCreatedEvent(tenantID string, orderID string) error
}

type OrderProcessedQueue interface {
	PublishOrderProcessedEvent(tenantID string, orderID string) error
}
