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
	webhooks2 "github.com/timhugh/digitalvenue/square/webhooks"
)

func NewLogger() zerolog.Logger {
	return log.With().Str("service", "square-events").Logger()
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
		webhooks2.NewHandlerProvider,
		webhooks2.NewPaymentCreatedService,
	)
	return handler{}, nil
}
