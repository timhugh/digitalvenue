package square_test

import (
	"context"
	"errors"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/squaretest"
	"github.com/timhugh/digitalvenue/util/test"
	"testing"
)

func TestEventGatherer_Gather_Success(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.When(squareAPI.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	customerCaptor := mock.Captor[*core.Customer]()
	mock.When(customerRepo.PutCustomer(customerCaptor.Capture())).ThenReturn(nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	orderCaptor := mock.Captor[*core.Order]()
	mock.When(orderRepo.PutOrder(orderCaptor.Capture())).ThenReturn(nil)

	orderCreatedQueue := mock.Mock[core.OrderCreatedQueue]()
	mock.When(orderCreatedQueue.PublishOrderCreatedEvent(test.TenantID, test.OrderID)).ThenReturn(nil)

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, orderRepo, customerRepo, squareAPI, orderCreatedQueue, orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
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

func TestEventGatherer_Gather_PaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(mock.AnyString(), mock.AnyString())).ThenReturn(&square.Payment{}, errors.New("failed to find payment"))

	squareAPI := mock.Mock[square.APIClient]()
	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, mock.Mock[square.MerchantRepository](), mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find payment")
}

func TestEventGatherer_Gather_MerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(&square.Merchant{}, errors.New("failed to find merchant"))

	squareAPI := mock.Mock[square.APIClient]()
	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), mock.Mock[square.APIClient](), mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_Gather_SquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(&square.Order{}, errors.New("failed to fetch order"))

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareAPI, mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_Gather_SquareFetchCustomerError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.When(squareAPI.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(&square.Customer{}, errors.New("failed to fetch customer"))

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), mock.Mock[core.CustomerRepository](), squareAPI, mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to fetch customer")
}

func TestEventGatherer_Gather_CustomerRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.When(squareAPI.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.When(customerRepo.PutCustomer(mock.Any[*core.Customer]())).ThenReturn(errors.New("failed to save customer"))

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, mock.Mock[core.OrderRepository](), customerRepo, squareAPI, mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save customer")
}

func TestEventGatherer_Gather_OrderRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.When(squareAPI.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.When(customerRepo.PutCustomer(mock.Any[*core.Customer]())).ThenReturn(nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	mock.When(orderRepo.PutOrder(mock.Any[*core.Order]())).ThenReturn(errors.New("failed to save order"))

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, orderRepo, customerRepo, squareAPI, mock.Mock[core.OrderCreatedQueue](), orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.Equal(err.Error(), "failed to save order")
}

func TestEventGatherer_Gather_OrderCreatedEventQueueError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	paymentRepo := mock.Mock[square.PaymentRepository]()
	mock.When(paymentRepo.GetSquarePayment(squaretest.SquareMerchantID, squaretest.SquarePaymentID)).ThenReturn(squaretest.NewSquarePayment(), nil)

	merchantRepo := mock.Mock[square.MerchantRepository]()
	mock.When(merchantRepo.GetSquareMerchant(squaretest.SquareMerchantID)).ThenReturn(squaretest.NewSquareMerchant(), nil)

	squareAPI := mock.Mock[square.APIClient]()
	mock.When(squareAPI.GetOrder(squaretest.SquareOrderID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareOrder(), nil)
	mock.When(squareAPI.GetCustomer(squaretest.SquareCustomerID, squaretest.SquareAPIToken)).ThenReturn(squaretest.NewSquareCustomer(), nil)

	customerRepo := mock.Mock[core.CustomerRepository]()
	mock.When(customerRepo.PutCustomer(mock.Any[*core.Customer]())).ThenReturn(nil)

	orderRepo := mock.Mock[core.OrderRepository]()
	mock.When(orderRepo.PutOrder(mock.Any[*core.Order]())).ThenReturn(nil)

	orderCreatedQueue := mock.Mock[core.OrderCreatedQueue]()
	mock.When(orderCreatedQueue.PublishOrderCreatedEvent(test.TenantID, squaretest.SquareOrderID)).ThenReturn(errors.New("failed to publish order created event"))

	orderBuilder := square.NewOrderBuilder(logger.Default(), squareAPI)

	gatherer := square.NewPaymentGatherer(paymentRepo, merchantRepo, orderRepo, customerRepo, squareAPI, orderCreatedQueue, orderBuilder)
	err := gatherer.Gather(context.TODO(), squaretest.SquareMerchantID, squaretest.SquarePaymentID)
	is.True(err != nil)
	is.Equal(err.Error(), "failed to publish order created event")
}
