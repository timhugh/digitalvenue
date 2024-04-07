package square_test

import (
	"errors"
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

	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[core.Customer]()
	mock.WhenDouble(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(test.CustomerID, nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	orderCaptor := mock.Captor[core.Order]()
	mock.WhenDouble(orderRepo.PutOrder(orderCaptor.Capture())).ThenReturn(test.OrderID, nil)

	orderCreatedQueue := mock.Mock[core.OrderCreatedQueue]()
	publishedOrderIDCaptor := mock.Captor[string]()
	mock.WhenSingle(orderCreatedQueue.Publish(publishedOrderIDCaptor.Capture())).ThenReturn(nil)

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi, orderCreatedQueue)
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

	expectedPublishedOrderID := test.OrderID
	publishedOrderID := publishedOrderIDCaptor.Last()
	is.Equal(expectedPublishedOrderID, publishedOrderID)
}

func TestEventGatherer_Gather_PaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.WhenDouble(paymentRepo.GetSquarePayment(squaretest.SquarePaymentID)).ThenReturn(square.Payment{}, errors.New("failed to find payment"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, mock.Mock[square.MerchantRepository](), mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), mock.Mock[core.OrderCreatedQueue]())
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

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), mock.Mock[core.OrderCreatedQueue]())
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

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi, mock.Mock[core.OrderCreatedQueue]())
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
	mock.WhenDouble(squareApi.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(square.Customer{}, errors.New("failed to fetch customer"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareApi, mock.Mock[core.OrderCreatedQueue]())
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
	mock.WhenDouble(customerRepo.PutCustomer(mock.Any[core.Customer]())).ThenReturn("", errors.New("failed to save customer"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), customerRepo, squareApi, mock.Mock[core.OrderCreatedQueue]())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save customer")
}

func TestEventGatherer_Gather_OrderRepoError(t *testing.T) {
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
	customerCaptor := mock.Captor[core.Customer]()
	mock.WhenDouble(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(test.CustomerID, nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	mock.WhenDouble(orderRepo.PutOrder(mock.Any[core.Order]())).ThenReturn("", errors.New("failed to save order"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi, mock.Mock[core.OrderCreatedQueue]())
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save order")
}

func TestEventGatherer_Gather_OrderCreatedQueueError(t *testing.T) {
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
	mock.WhenDouble(customerRepo.PutCustomer(mock.Any[core.Customer]())).ThenReturn(test.CustomerID, nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	mock.WhenDouble(orderRepo.PutOrder(mock.Any[core.Order]())).ThenReturn(test.OrderID, nil)

	orderCreatedQueue := mock.Mock[core.OrderCreatedQueue]()
	mock.WhenSingle(orderCreatedQueue.Publish(mock.Any[string]())).ThenReturn(errors.New("failed to publish order created event"))

	gatherer := square.NewPaymentGatherer(zerolog.Logger{}, paymentRepo, merchantRepo, orderRepo, customerRepo, squareApi, orderCreatedQueue)
	err := gatherer.Gather(squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to publish order created event")
}
