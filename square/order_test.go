package square

import (
	"bytes"
	"github.com/go-test/deep"
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/core"
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
		SquareLocationID: "location_id",
		OrderItems: []OrderItem{
			{
				ItemID:   "item_uid_1",
				Name:     "Item 1",
				Quantity: "1",
			},
			{
				ItemID:   "item_uid_2",
				Name:     "Item 2",
				Quantity: "2",
			},
		},
	}
	if diff := deep.Equal(order, expectedOrder); diff != nil {
		t.Error(diff)
	}
}

func TestMapOrder(t *testing.T) {
	is := is.New(t)

	squareOrder := Order{
		SquareOrderID:    "order_id",
		SquareCustomerID: "customer_id",
		SquareLocationID: "location_id",
		OrderItems: []OrderItem{
			{
				ItemID:   "item_uid_1",
				Name:     "Item 1",
				Quantity: "1",
			},
			{
				ItemID:   "item_uid_2",
				Name:     "Item 2",
				Quantity: "2",
			},
		},
	}

	order, err := MapOrder(squareOrder)
	is.NoErr(err)

	expectedOrder := core.Order{
		Items: []core.OrderItem{
			{
				Name: "Item 1",
				Meta: core.OrderItemMeta{
					SquareItemID: "item_uid_1",
				},
			},
			{
				Name: "Item 2",
				Meta: core.OrderItemMeta{
					SquareItemID: "item_uid_2",
				},
			},
			{
				Name: "Item 2",
				Meta: core.OrderItemMeta{
					SquareItemID: "item_uid_2",
				},
			},
		},
		Meta: core.OrderMeta{
			SquareOrderID:    "order_id",
			SquareCustomerID: "customer_id",
		},
	}
	if diff := deep.Equal(order, expectedOrder); diff != nil {
		t.Error(diff)
	}
}
