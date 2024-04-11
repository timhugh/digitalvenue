//go:build wireinject
// +build wireinject

package main

import (
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/dv_aws"
	"github.com/timhugh/digitalvenue/dv_aws/dv_dynamodb"
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

		dv_aws.DefaultConfig,

		dv_dynamodb.NewClient,
		wire.Bind(new(dv_dynamodb.Client), new(*awsdynamodb.Client)),
		dv_dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dv_dynamodb.Repository)),

		squareapi.NewClient,
		wire.Bind(new(square.APIClient), new(*squareapi.Client)),

		NewSquareEventGathererHandler,
	)
	return SquareEventGathererHandler{}, nil
}
