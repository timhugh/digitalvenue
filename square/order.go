package square

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
