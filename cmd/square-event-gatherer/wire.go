//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/aws"
	"github.com/timhugh/digitalvenue/aws/dynamodb"
	square2 "github.com/timhugh/digitalvenue/aws/dynamodb/square"
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
		square2.NewPaymentsRepositoryConfig,
		square2.NewPaymentsRepository,
		square.NewEventGatherer,
	)
	return handler{}, nil
}
