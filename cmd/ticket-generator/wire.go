//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newLogger() zerolog.Logger {
	return log.With().Str("service", "ticket-generator").Logger()
}

func initializeHandler() (*TicketGeneratorHandler, error) {
	wire.Build(
		newLogger,
		NewTicketGeneratorHandler,
	)
	return &TicketGeneratorHandler{}, nil
}
