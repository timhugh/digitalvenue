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
		dynamodb.NewSquareMerchantsRepository,
		dynamodb.NewSquareMerchantsRepositoryConfig,
		dynamodb.NewSquarePaymentsRepository,
		dynamodb.NewSquarePaymentsRepositoryConfig,
		newEventServiceConfig,
		newHandler,
		sqs.NewClient,
		sqs.NewSquarePaymentCreatedQueue,
		sqs.NewSquarePaymentCreatedQueueConfig,
		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,
	)
	return handler{}, nil
}
