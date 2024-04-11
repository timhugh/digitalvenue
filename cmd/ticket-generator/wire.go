//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
)

func initializeHandler(logger zerolog.Logger) (*TicketGeneratorHandler, error) {
	wire.Build(
		NewTicketGeneratorHandler,
	)
	return &TicketGeneratorHandler{}, nil
}
