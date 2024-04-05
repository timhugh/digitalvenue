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
	Source           string
	SquareCustomerID string
}

type CustomerRepository interface {
	Put(customer Customer) (string, error)
}
