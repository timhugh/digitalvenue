//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/aws"
	"github.com/timhugh/digitalvenue/aws/dynamodb"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/handlers"
	"github.com/timhugh/digitalvenue/square/squareapi"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-event-gatherer").Logger()
}

func initializeHandler() (handlers.SquareEventGathererHandler, error) {
	wire.Build(
		newLogger,
		square.NewEventGatherer,

		aws.DefaultConfig,

		dynamodb.NewClient,
		dynamodb.NewSquareMerchantRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.SquareMerchantRepository)),
		dynamodb.NewSquarePaymentRepository,
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.SquarePaymentRepository)),
		dynamodb.NewOrderRepository,
		wire.Bind(new(core.OrderRepository), new(*dynamodb.OrderRepository)),
		dynamodb.NewCustomerRepository,
		wire.Bind(new(core.CustomerRepository), new(*dynamodb.CustomerRepository)),

		squareapi.NewClient,
		wire.Bind(new(square.APIClient), new(*squareapi.Client)),
		square.NewOrderMapper,

		handlers.NewSquareEventGathererHandler,
	)
	return handlers.SquareEventGathererHandler{}, nil
}
