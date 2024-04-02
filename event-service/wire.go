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
	"github.com/timhugh/digitalvenue/services/webhooks"
)

func NewLogger() zerolog.Logger {
	return log.With().Str("service", "event-service").Logger()
}

func initializeHandler() (handler, error) {
	wire.Build(
		NewLogger,
		aws.NewConfig,
		dynamodb.NewClient,
		dynamodb.NewMerchantsRepository,
		dynamodb.NewMerchantsRepositoryConfig,
		dynamodb.NewPaymentsRepository,
		dynamodb.NewPaymentsRepositoryConfig,
		newEventServiceConfig,
		newHandler,
		sqs.NewClient,
		sqs.NewPaymentCreatedQueue,
		sqs.NewPaymentCreatedQueueConfig,
		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedService,
	)
	return handler{}, nil
}
