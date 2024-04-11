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
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-events").Logger()
}

func initializeHandler() (*SquareEventsHandler, error) {
	wire.Build(
		newLogger,

		aws.DefaultConfig,

		dynamodb.NewClient,
		wire.Bind(new(dynamodb.Client), new(*awsdynamodb.Client)),
		dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.Repository)),

		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,

		NewSquareEventsHandler,
	)
	return &SquareEventsHandler{}, nil
}
