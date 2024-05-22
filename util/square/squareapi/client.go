package squareapi

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/square"
	"io"
	"net/http"
)

const (
	maxBodyLength = 1048576
	bearerToken   = "Bearer %s"

	squareApiBaseUrl            = "https://connect.squareup.com"
	getOrderRouteFormat         = "/v2/orders/%s"         // squareOrderID
	getCustomerRouteFormat      = "/v2/customers/%s"      // squareCustomerID
	getCatalogObjectRouteFormat = "/v2/catalog/object/%s" // catalogObjectID
)

type ApiError struct {
	Errors []struct {
		Code     string `json:"code"`
		Detail   string `json:"detail"`
		Category string `json:"category"`
	} `json:"errors"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseUrl       string
	maxBodyLength int64
	httpClient    httpClient
}

func NewClient() square.APIClient {
	return &Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    http.DefaultClient,
	}
}

func (client *Client) fetchJson(path string, apiToken string, target interface{}) error {
	body, err := client.fetchBody(path, apiToken)
	if err != nil {
		return errors.Wrap(err, "failed to get data from square API")
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal http response JSON")
	}

	return nil
}

func (client *Client) fetchBody(path string, apiToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Set("Authorization", fmt.Sprintf(bearerToken, apiToken))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing http request")
	}

	body, err := client.readBody(resp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http response body")
	}

	if resp.StatusCode != http.StatusOK {
		var apiError ApiError
		err = json.Unmarshal(body, &apiError)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal API error JSON")
		}
		// TODO: this is just the first error; should probably return all of them if there are more
		return nil, errors.Errorf("API error: %s", apiError.Errors[0].Detail)
	}

	return body, nil
}

func (client *Client) readBody(resp *http.Response) ([]byte, error) {
	buf, err := io.ReadAll(io.LimitReader(resp.Body, client.maxBodyLength))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (client *Client) GetCustomer(squareCustomerID string, apiToken string) (*square.Customer, error) {
	path := client.baseUrl + fmt.Sprintf(getCustomerRouteFormat, squareCustomerID)

	var customerContainer struct {
		Customer square.Customer `json:"customer"`
	}
	err := client.fetchJson(path, apiToken, &customerContainer)
	if err != nil {
		return nil, err
	}

	return &customerContainer.Customer, nil
}

func (client *Client) GetOrder(squareOrderID string, apiToken string) (*square.Order, error) {
	path := client.baseUrl + fmt.Sprintf(getOrderRouteFormat, squareOrderID)

	var orderContainer struct {
		Order square.Order `json:"order"`
	}
	err := client.fetchJson(path, apiToken, &orderContainer)
	if err != nil {
		return nil, err
	}

	return &orderContainer.Order, nil
}

func (client *Client) GetCatalogObject(objectID string, apiToken string, out interface{}) error {
	path := client.baseUrl + fmt.Sprintf(getCatalogObjectRouteFormat, objectID)

	var objectContainer struct {
		Object square.CatalogObject `json:"object"`
	}
	body, err := client.fetchBody(path, apiToken)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &objectContainer)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal http response JSON as CatalogObject")
	}

	switch objectContainer.Object.Type {
	case square.CatalogItemVariationType:
		var itemVariationContainer struct {
			Object square.CatalogItemVariation `json:"object"`
		}
		err := json.Unmarshal(body, &itemVariationContainer)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal http response JSON as CatalogItemVariation")
		}
		out = itemVariationContainer.Object
		return nil
	case square.CatalogItemType:
		var itemContainer struct {
			Object square.CatalogItem `json:"object"`
		}
		err := json.Unmarshal(body, &itemContainer)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal http response JSON as CatalogItem")
		}
	}

	return errors.New("retrieved unknown catalog object type")
}
