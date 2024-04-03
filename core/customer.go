package core

type Customer struct {
	CustomerID string
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
	Create(customer Customer) error
}
