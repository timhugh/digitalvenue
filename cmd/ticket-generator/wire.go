//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core/services"
)

func initializeHandler(logger zerolog.Logger) (*TicketGeneratorHandler, error) {
	wire.Build(
		services.NewTicketGenerator,
		NewTicketGeneratorHandler,
	)
	return &TicketGeneratorHandler{}, nil
}
