//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/aws"
	"github.com/timhugh/digitalvenue/aws/dynamodb"
	dynamoSquare "github.com/timhugh/digitalvenue/aws/dynamodb/square"
	"github.com/timhugh/digitalvenue/square"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-event-gatherer").Logger()
}

func initializeHandler() (handler, error) {
	wire.Build(
		newLogger,
		newHandler,
		aws.NewConfig,
		dynamodb.NewClient,
		square.NewHttpClient,
		square.NewClientConfig,
		square.NewClient,
		dynamoSquare.NewPaymentsRepositoryConfig,
		dynamoSquare.NewPaymentsRepository,
		dynamoSquare.NewOrderRepositoryConfig,
		dynamoSquare.NewOrderRepository,
		dynamoSquare.NewMerchantsRepositoryConfig,
		dynamoSquare.NewMerchantsRepository,
		square.NewEventGatherer,
	)
	return handler{}, nil
}
