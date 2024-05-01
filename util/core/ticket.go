package core

type Ticket struct {
	ID         string
	TenantID   string
	OrderID    string
	CustomerID string
	Name       string
	QRCodeURL  string
}

type TicketRepository interface {
	PutTickets(tickets []*Ticket) error
	GetTickets(tenantID string, orderID string) ([]*Ticket, error)
}
