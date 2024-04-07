package square_test

import (
	"errors"
	"fmt"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squaretest"
	"github.com/timhugh/digitalvenue/test"
	"testing"
)

func TestEventGatherer_Gather_Success(t *testing.T) {
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

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi)
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.NoErr(err)

	expectedOrder := squaretest.NewOrder()
	expectedOrder.CustomerID = test.CustomerID
	persistedOrder := orderCaptor.Last()
	err = test.Diff(expectedOrder, persistedOrder)
	is.NoErr(err)

	expectedCustomer := squaretest.NewCustomer()
	persistedCustomer := customerCaptor.Last()
	err = test.Diff(expectedCustomer, persistedCustomer)
	is.NoErr(err)
}

func TestEventGatherer_Gather_PaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(square.Payment{}, errors.New("failed to find payment"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, mock.Mock[square.MerchantRepository](), mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient]())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find payment")
}

func TestEventGatherer_Gather_MerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(square.Merchant{}, errors.New("failed to find merchant"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient]())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_Gather_SquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(square.Order{}, errors.New("failed to fetch order"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi)
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_Gather_SquareFetchCustomerError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.WhenDouble(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareApi := mock.Mock[square.APIClient]()
	mock.WhenDouble(squareApi.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(square.Customer{}, fmt.Errorf("failed to fetch customer"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi)
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch customer")
}

func TestEventGatherer_Gather_CustomerRepoError(t *testing.T) {
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

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), customerRepo, squareApi)
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save customer")
}

func TestEventGatherer_Gather_OrderRepoError(t *testing.T) {}
