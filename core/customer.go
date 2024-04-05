package core

type Customer struct {
	CustomerID string
	TenantID   string
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Meta       CustomerMeta
}

type CustomerMeta struct {
	SquareCustomerID string
}

type CustomerRepository interface {
	PutCustomer(customer Customer) (string, error)
}
