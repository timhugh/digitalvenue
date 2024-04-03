package square

import (
	"bytes"
	"github.com/matryer/is"
	"io"
	"net/http"
	"testing"
)

func TestClient_GetOrder_Success(t *testing.T) {
	is := is.New(t)

	httpClient := NewTestClient(func(r *http.Request) *http.Response {
		is.Equal(r.Method, http.MethodGet)
		is.Equal(r.Header.Get("Authorization"), "Bearer api_token")
		is.Equal(r.URL.Path, "/v2/orders/order_id")

		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(orderJson)),
		}
	})

	squareClient := NewClient(NewClientConfig(), httpClient)

	order, err := squareClient.GetOrder("order_id", "api_token")
	is.NoErr(err)

	expectedOrder := Order{
		SquareOrderID:    "order_id",
		SquareCustomerID: "customer_id",
	}
	is.Equal(order, expectedOrder)
}
