//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_sqs"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/webhooks"
)

func initializeHandler(log *logger.ContextLogger) (*SquareEventsHandler, error) {
	wire.Build(
		dv_aws.DefaultConfig,

		dv_dynamodb.NewClient,
		dv_dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dv_dynamodb.Repository)),

		dv_sqs.NewClient,
		dv_sqs.NewSquarePaymentCreatedQueue,
		wire.Bind(new(square.PaymentCreatedQueue), new(*dv_sqs.SquarePaymentCreatedQueue)),

		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,

		NewSquareEventsHandler,
	)
	return &SquareEventsHandler{}, nil
}
