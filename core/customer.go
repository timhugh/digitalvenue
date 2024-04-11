package core

type Customer struct {
	ID       string
	TenantID string
	Name     string
	Email    string
	Phone    string
	Meta     map[string]string
}

type CustomerRepository interface {
	PutCustomer(customer Customer) (string, error)
}
