package square

import (
	"fmt"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"testing"
)

func TestEventGatherer_GatherSuccess(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(nil)

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.NoErr(err)
}

func TestEventGatherer_GatherPaymentRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(squarePaymentID)).ThenReturn(Payment{}, fmt.Errorf("failed to find payment"))

	merchantRepo := mock.Mock[MerchantsRepository]()
	squareApi := mock.Mock[Client]()
	orderRepo := mock.Mock[OrdersRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to find payment")
}

func TestEventGatherer_GatherMerchantRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(Merchant{}, fmt.Errorf("failed to find merchant"))

	squareApi := mock.Mock[Client]()
	orderRepo := mock.Mock[OrdersRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to find merchant")
}

func TestEventGatherer_GatherSquareFetchOrderError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	squarePaymentID := "squarePaymentID"

	log := zerolog.Logger{}

	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(Order{}, fmt.Errorf("failed to fetch order"))

	orderRepo := mock.Mock[OrdersRepository]()

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, squareApi)

	err := gatherer.Gather(squarePaymentID)
	is.Equal(err.Error(), "failed to fetch order")
}

func TestEventGatherer_GatherOrderRepoError(t *testing.T) {
	is := is.New(t)
	mock.SetUp(t)

	log := zerolog.Logger{}
	payment := Payment{
		SquareMerchantID: "squareMerchantID",
		SquareOrderID:    "squareOrderID",
		SquarePaymentID:  "squarePaymentID",
	}
	paymentRepo := mock.Mock[PaymentsRepository]()
	mock.WhenDouble(paymentRepo.FindByID(payment.SquarePaymentID)).ThenReturn(payment, nil)

	merchant := Merchant{
		SquareMerchantID: "squareMerchantID",
		SquareAPIToken:   "squareAPIToken",
	}
	merchantRepo := mock.Mock[MerchantsRepository]()
	mock.WhenDouble(merchantRepo.FindByID(payment.SquareMerchantID)).ThenReturn(merchant, nil)

	expectedOrder := Order{
		SquareOrderID:    "squareOrderID",
		SquareCustomerID: "squareCustomerID",
	}
	squareApi := mock.Mock[Client]()
	mock.WhenDouble(squareApi.GetOrder(payment.SquareOrderID, merchant.SquareAPIToken)).ThenReturn(expectedOrder, nil)
	orderRepo := mock.Mock[OrdersRepository]()
	mock.WhenSingle(orderRepo.Create(mock.Equal(expectedOrder))).ThenReturn(fmt.Errorf("failed to create order"))

	gatherer := NewEventGatherer(log, paymentRepo, merchantRepo, orderRepo, squareApi)

	err := gatherer.Gather("squarePaymentID")
	is.Equal(err.Error(), "failed to create order")
}
