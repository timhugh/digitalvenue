package squareapi

import (
	"bytes"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"io"
	"net/http"
	"os"
	"testing"
)

var OrderRawJSON, _ = os.ReadFile("test-order-response.json")
var OrderJSON = string(OrderRawJSON)

var CustomerRawJSON, _ = os.ReadFile("test-customer-response.json")
var CustomerJSON = string(CustomerRawJSON)

var ErrorNotFoundRawJSON, _ = os.ReadFile("test-error-not-found.json")
var ErrorNotFoundJSON = string(ErrorNotFoundRawJSON)

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
			Body:       io.NopCloser(bytes.NewBufferString(CustomerJSON)),
			StatusCode: http.StatusOK,
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
	is.Equal(expectedCustomer, customer)
}

func TestClient_GetOrder_Success(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.Header.Get("Authorization"), "Bearer api_token")
		is.Equal(r.URL.Path, "/v2/orders/squareOrderID")

		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString(OrderJSON)),
			StatusCode: http.StatusOK,
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
	err = test.Diff(expectedOrder, order)
	is.NoErr(err)
}

func TestClient_GetOrder_ErrorNotFound(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString(ErrorNotFoundJSON)),
			StatusCode: http.StatusNotFound,
		}
	})

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    httpClient,
	}

	_, err := squareClient.GetOrder("squareOrderID", "api_token")
	is.Equal(err.Error(), "API error: Order not found for id some_order_id")
}

func TestClient_GetOrder_ErrorNetworkFailure(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	thrownError := errors.New("timeout")
	mockHttpClient := mock.Mock[httpClient]()
	mock.WhenDouble(mockHttpClient.Do(mock.Any[*http.Request]())).ThenReturn(&http.Response{}, thrownError)

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    mockHttpClient,
	}

	_, err := squareClient.GetOrder("squareOrderID", "api_token")
	is.True(errors.Is(err, thrownError))
}

func TestClient_GetOrder_ErrorNoBody(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		return &http.Response{
			Body:       nil,
			StatusCode: http.StatusOK,
		}
	})

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    httpClient,
	}

	_, err := squareClient.GetOrder("squareOrderID", "api_token")
	is.True(err != nil)
}

func TestClient_GetOrder_ErrorInvalidBody(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		return &http.Response{
			Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
			StatusCode: http.StatusOK,
		}
	})

	squareClient := Client{
		baseUrl:       squareApiBaseUrl,
		maxBodyLength: maxBodyLength,
		httpClient:    httpClient,
	}

	_, err := squareClient.GetOrder("squareOrderID", "api_token")
	is.True(err != nil)
}
