//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_sqs"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/squareapi"
)

func initializeHandler(log *logger.ContextLogger) (SquareEventGathererHandler, error) {
	wire.Build(
		square.NewPaymentGatherer,

		dv_aws.DefaultConfig,

		dv_dynamodb.NewClient,
		dv_dynamodb.NewRepository,
		wire.Bind(new(square.PaymentRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(square.MerchantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dv_dynamodb.Repository)),
		dv_sqs.NewClient,
		dv_sqs.NewOrderCreatedQueue,
		wire.Bind(new(core.OrderCreatedQueue), new(*dv_sqs.OrderCreatedQueue)),

		squareapi.NewClient,
		wire.Bind(new(square.APIClient), new(*squareapi.Client)),

		NewSquareEventGathererHandler,
	)
	return SquareEventGathererHandler{}, nil
}
