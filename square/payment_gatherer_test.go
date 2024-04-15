package square_test

import (
	"context"
	"errors"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func TestEventGatherer_Gather_Success(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[*core.Customer]()
	mock.WhenSingle(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	orderCaptor := mock.Captor[*core.Order]()
	mock.WhenSingle(orderRepo.PutOrder(orderCaptor.Capture())).ThenReturn(nil)

	gatherer := square.NewPaymentGatherer(merchantRepo, orderRepo, customerRepo, squareApi)
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.NoErr(err)

	expectedOrder := squaretest.NewOrder()
	persistedOrder := orderCaptor.Last()
	err = test.Diff(expectedOrder, persistedOrder)
	is.NoErr(err)

	expectedCustomer := squaretest.NewCustomer()
	persistedCustomer := customerCaptor.Last()
	err = test.Diff(expectedCustomer, persistedCustomer)
	is.NoErr(err)
}

func TestEventGatherer_Gather_MerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(&square.Merchant{}, errors.New("failed to find merchant"))

	gatherer := square.NewPaymentGatherer(merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient]())
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_Gather_SquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(&square.Order{}, errors.New("failed to fetch order"))

	gatherer := square.NewPaymentGatherer(merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi)
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_Gather_SquareFetchCustomerError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(&square.Customer{}, errors.New("failed to fetch customer"))

	gatherer := square.NewPaymentGatherer(merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi)
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.Equal(err.Error(), "failed to fetch customer")
}

func TestEventGatherer_Gather_CustomerRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.WhenSingle(customerRepo.PutCustomer(mock.Any[*core.Customer]())).ThenReturn(errors.New("failed to save customer"))

	gatherer := square.NewPaymentGatherer(merchantRepo, mock.Mock[core.OrderRepository](), customerRepo, squareApi)
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.Equal(err.Error(), "failed to save customer")
}

func TestEventGatherer_Gather_OrderRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[*core.Customer]()
	mock.WhenSingle(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	mock.WhenSingle(orderRepo.PutOrder(mock.Any[*core.Order]())).ThenReturn(errors.New("failed to save order"))

	gatherer := square.NewPaymentGatherer(merchantRepo, orderRepo, customerRepo, squareApi)
	err := gatherer.Gather(context.TODO(), squaretest.NewSquarePayment())
	is.Equal(err.Error(), "failed to save order")
}
