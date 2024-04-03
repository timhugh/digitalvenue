package square

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	square        = "square"
	maxBodyLength = 1048576
	bearerToken   = "Bearer %s"

	squareApiBaseUrl       = "https://connect.squareup.com"
	getOrderRouteFormat    = "/v2/orders/%s"    // squareOrderID
	getCustomerRouteFormat = "/v2/customers/%s" // squareCustomerID
)

type Client interface {
	GetOrder(orderId string, apiToken string) (Order, error)
	GetCustomer(customerId string, apiToken string) (Customer, error)
}

type ClientConfig struct {
	BaseUrl       string
	MaxBodyLength int64
}

func NewClientConfig() ClientConfig {
	return ClientConfig{
		BaseUrl:       squareApiBaseUrl,
		MaxBodyLength: maxBodyLength,
	}
}

func NewHttpClient() *http.Client {
	return http.DefaultClient
}

type client struct {
	baseUrl       string
	maxBodyLength int64
	httpClient    *http.Client
}

func NewClient(config ClientConfig, httpClient *http.Client) Client {
	return client{
		baseUrl:       config.BaseUrl,
		maxBodyLength: config.MaxBodyLength,
		httpClient:    httpClient,
	}
}

func (client client) fetchJson(path string, apiToken string, target interface{}) error {
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

func (client client) readBody(resp *http.Response) ([]byte, error) {
	buf, err := io.ReadAll(io.LimitReader(resp.Body, client.maxBodyLength))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return buf, nil
}
