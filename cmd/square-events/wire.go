//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/aws"
	"github.com/timhugh/digitalvenue/aws/dynamodb"
	"github.com/timhugh/digitalvenue/aws/sqs"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-events").Logger()
}

func initializeHandler() (handler, error) {
	wire.Build(
		newLogger,

		aws.NewConfig,

		dynamodb.NewClient,
		dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.Repository)),

		sqs.NewClient,
		sqs.NewQueue,
		wire.Bind(new(square.PaymentCreatedQueue), new(*sqs.Queue)),

		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,

		newHandler,
	)
	return handler{}, nil
}
