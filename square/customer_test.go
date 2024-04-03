package square

import (
	"bytes"
	"github.com/matryer/is"
	"io"
	"net/http"
	"testing"
)

func TestClient_GetCustomer_Success(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.Header.Get("Authorization"), "Bearer api_token")
		is.Equal(r.URL.Path, "/v2/customers/customer_id")

		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(customerJson)),
		}
	})

	squareClient := NewClient(NewClientConfig(), httpClient)

	customer, err := squareClient.GetCustomer("customer_id", "api_token")
	is.NoErr(err)

	expectedCustomer := Customer{
		SquareCustomerID: "customer_id",
		FirstName:        "Tim",
		LastName:         "Heuett",
		Email:            "info@timheuett.com",
		Phone:            "+12062062062",
	}
	is.Equal(customer, expectedCustomer)
}
