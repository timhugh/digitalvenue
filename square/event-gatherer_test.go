package square_test

import (
	"fmt"
	"github.com/go-test/deep"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func TestEventGatherer_GatherSuccess(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	orderCaptor := mock.Captor[core.Order]()
	mock.WhenDouble(orderRepo.PutOrder(orderCaptor.Capture())).ThenReturn("order_id", nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[core.Customer]()
	mock.WhenDouble(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(test.CustomerID, nil)

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi, square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.NoErr(err)

	expectedOrder := squaretest.NewOrderGathered()
	persistedOrder := orderCaptor.Last()
	if diff := deep.Equal(expectedOrder, persistedOrder); diff != nil {
		t.Error(diff)
	}

	expectedCustomer := squaretest.NewCustomer()
	persistedCustomer := customerCaptor.Last()
	if diff := deep.Equal(expectedCustomer, persistedCustomer); diff != nil {
		t.Error(diff)
	}
}

func TestEventGatherer_GatherPaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(square.Payment{}, fmt.Errorf("failed to find payment"))

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, mock.Mock[square.MerchantRepository](), mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find payment")
}

func TestEventGatherer_GatherMerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(square.Merchant{}, fmt.Errorf("failed to find merchant"))

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_GatherSquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(square.Order{}, fmt.Errorf("failed to fetch order"))

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi, square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_GatherSquareFetchCustomerError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(square.Customer{}, fmt.Errorf("failed to fetch customer"))

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi, square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch customer")
}

func TestEventGatherer_GatherCustomerRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.WhenDouble(customerRepo.PutCustomer(mock.Any[core.Customer]())).ThenReturn("", fmt.Errorf("failed to save customer"))

	gatherer := square.NewEventGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), customerRepo, squareApi, square.NewOrderMapper())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save customer")
}
