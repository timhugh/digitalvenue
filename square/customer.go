package square

import (
	"fmt"
	"github.com/timhugh/digitalvenue/core"
)

type Customer struct {
	SquareCustomerID string
	FirstName        string
	LastName         string
	Email            string
	Phone            string
}

type customerContainer struct {
	Customer struct {
		ID           string `json:"id"`
		GivenName    string `json:"given_name"`
		FamilyName   string `json:"family_name"`
		EmailAddress string `json:"email_address"`
		PhoneNumber  string `json:"phone_number"`
	} `json:"customer"`
}

func (client client) GetCustomer(squareCustomerID string, apiToken string) (Customer, error) {
	path := client.baseUrl + fmt.Sprintf(getCustomerRouteFormat, squareCustomerID)

	var customerContainer customerContainer
	err := client.fetchJson(path, apiToken, &customerContainer)
	if err != nil {
		return Customer{}, err
	}

	return Customer{
		SquareCustomerID: customerContainer.Customer.ID,
		FirstName:        customerContainer.Customer.GivenName,
		LastName:         customerContainer.Customer.FamilyName,
		Email:            customerContainer.Customer.EmailAddress,
		Phone:            customerContainer.Customer.PhoneNumber,
	}, nil
}

func MapCustomer(customer Customer) core.Customer {
	return core.Customer{
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Meta: core.CustomerMeta{
			Source:           square,
			SquareCustomerID: customer.SquareCustomerID,
		},
	}
}
