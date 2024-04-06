package square

import "github.com/timhugh/digitalvenue/core"

const (
	CustomerIDKey = "SquareCustomerID"
)

type Customer struct {
	SquareCustomerID string
	FirstName        string
	LastName         string
	Email            string
	Phone            string
}

func MapCustomer(customer Customer) core.Customer {
	return core.Customer{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Meta: map[string]string{
			CustomerIDKey: customer.SquareCustomerID,
		},
	}
}
