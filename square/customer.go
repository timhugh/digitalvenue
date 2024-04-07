package square

import "github.com/timhugh/digitalvenue/core"

const (
	CustomerIDKey = "SquareCustomerID"
)

type Customer struct {
	SquareCustomerID string `json:"id"`
	FirstName        string `json:"given_name"`
	LastName         string `json:"family_name"`
	Email            string `json:"email_address"`
	Phone            string `json:"phone_number"`
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
