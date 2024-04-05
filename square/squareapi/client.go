package squareapi

import (
	"encoding/json"
	"fmt"
	"github.com/timhugh/digitalvenue/square"
	"io"
	"net/http"
)

const (
	maxBodyLength = 1048576
	bearerToken   = "Bearer %s"

	squareApiBaseUrl       = "https://connect.squareup.com"
	getOrderRouteFormat    = "/v2/orders/%s"    // squareOrderID
	getCustomerRouteFormat = "/v2/customers/%s" // squareCustomerID
)

type Client struct {
	baseUrl       string
	maxBodyLength int64
	httpClient    *http.Client
}

func NewClient() *Client {
	return &Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    http.DefaultClient,
	}
}

func (client *Client) fetchJson(path string, apiToken string, target interface{}) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf(bearerToken, apiToken))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	body, err := client.readBody(resp)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) readBody(resp *http.Response) ([]byte, error) {
	buf, err := io.ReadAll(io.LimitReader(resp.Body, client.maxBodyLength))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return buf, nil
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

func (client *Client) GetCustomer(squareCustomerID string, apiToken string) (square.Customer, error) {
	path := client.baseUrl + fmt.Sprintf(getCustomerRouteFormat, squareCustomerID)

	var customerContainer customerContainer
	err := client.fetchJson(path, apiToken, &customerContainer)
	if err != nil {
		return square.Customer{}, err
	}

	return square.Customer{
		SquareCustomerID: customerContainer.Customer.ID,
		FirstName:        customerContainer.Customer.GivenName,
		LastName:         customerContainer.Customer.FamilyName,
		Email:            customerContainer.Customer.EmailAddress,
		Phone:            customerContainer.Customer.PhoneNumber,
	}, nil
}

type orderContainer struct {
	Order square.Order `json:"order"`
}

func (client *Client) GetOrder(squareOrderID string, apiToken string) (square.Order, error) {
	path := client.baseUrl + fmt.Sprintf(getOrderRouteFormat, squareOrderID)

	var orderContainer orderContainer
	err := client.fetchJson(path, apiToken, &orderContainer)
	if err != nil {
		return square.Order{}, err
	}

	return orderContainer.Order, nil
}
