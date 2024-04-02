//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "square-event-gatherer").Logger()
}

func initializeHandler() (handler, error) {
	wire.Build(
		newLogger,
		newHandler,
	)
	return handler{}, nil
}
