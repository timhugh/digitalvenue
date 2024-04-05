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
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-event-gatherer").Logger()
}

func initializeHandler() (handler, error) {
	wire.Build(
		newLogger,
		square.NewEventGatherer,

		aws.NewConfig,

		dynamodb.NewClient,
		dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dynamodb.Repository)),

		square.NewClient,
		square.NewOrderMapper,

		newHandler,
	)
	return handler{}, nil
}
