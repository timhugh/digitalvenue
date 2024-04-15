package square

import (
	"github.com/timhugh/digitalvenue/util/core"
)

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

func MapCustomer(customer *Customer, tenantID string) *core.Customer {
	return &core.Customer{
		ID:       customer.SquareCustomerID, // We use the SquareCustomerID as the ID for customers that come from Square to make duplicate detection easier
		TenantID: tenantID,
		Name:     customer.FirstName + " " + customer.LastName,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Meta: map[string]string{
			CustomerIDKey: customer.SquareCustomerID,
		},
	}
}
