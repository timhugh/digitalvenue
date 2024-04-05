package squareapi

import (
	"bytes"
	"github.com/go-test/deep"
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"io"
	"net/http"
	"os"
	"testing"
)

var OrderRawJson, _ = os.ReadFile("test-order-response.json")
var OrderJson = string(OrderRawJson)

var CustomerRawJson, _ = os.ReadFile("test-customer-response.json")
var CustomerJson = string(CustomerRawJson)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClient_GetCustomer_Success(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.Header.Get("Authorization"), "Bearer api_token")
		is.Equal(r.URL.Path, "/v2/customers/squareCustomerID")

		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(CustomerJson)),
		}
	})

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    httpClient,
	}

	customer, err := squareClient.GetCustomer("squareCustomerID", "api_token")
	is.NoErr(err)

	expectedCustomer := squaretest.NewSquareCustomer()
	is.Equal(customer, expectedCustomer)
}

func TestClient_GetOrder_Success(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.Header.Get("Authorization"), "Bearer api_token")
		is.Equal(r.URL.Path, "/v2/orders/squareOrderID")

		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(OrderJson)),
		}
	})

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    httpClient,
	}

	order, err := squareClient.GetOrder("squareOrderID", "api_token")
	is.NoErr(err)

	expectedOrder := squaretest.NewSquareOrder()
	if diff := deep.Equal(order, expectedOrder); diff != nil {
		t.Error(diff)
	}
}
