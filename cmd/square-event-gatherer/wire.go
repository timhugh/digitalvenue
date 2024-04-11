//go:build wireinject
// +build wireinject

package main

import (
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/aws"
	"github.com/timhugh/digitalvenue/aws/dynamodb"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squareapi"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-event-gatherer").Logger()
}

func initializeHandler() (SquareEventGathererHandler, error) {
	wire.Build(
		newLogger,
		square.NewPaymentGatherer,

		aws.DefaultConfig,

		dynamodb.NewClient,
		wire.Bind(new(dynamodb.Client), new(*awsdynamodb.Client)),
		dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dynamodb.Repository)),

		squareapi.NewClient,
		wire.Bind(new(square.APIClient), new(*squareapi.Client)),

		NewSquareEventGathererHandler,
	)
	return SquareEventGathererHandler{}, nil
}
