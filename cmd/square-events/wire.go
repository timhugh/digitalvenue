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
	"github.com/timhugh/digitalvenue/square/handlers"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-events").Logger()
}

func initializeHandler() (handlers.SquareEventsHandler, error) {
	wire.Build(
		newLogger,

		aws.DefaultConfig,

		dynamodb.NewClient,
		dynamodb.NewSquareMerchantRepository,
		wire.Bind(new(square.MerchantRepository), new(*dynamodb.SquareMerchantRepository)),
		dynamodb.NewSquarePaymentRepository,
		wire.Bind(new(square.PaymentRepository), new(*dynamodb.SquarePaymentRepository)),

		sqs.NewClient,
		sqs.NewSquarePaymentCreatedQueue,
		wire.Bind(new(square.PaymentCreatedQueue), new(*sqs.SquarePaymentCreatedQueue)),

		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,

		handlers.NewSquareEventsHandler,
	)
	return handlers.SquareEventsHandler{}, nil
}
