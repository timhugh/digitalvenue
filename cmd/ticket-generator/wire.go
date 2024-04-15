//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/core/services"
	"github.com/timhugh/digitalvenue/logger"
)

func initializeHandler(logger *logger.ContextLogger) (*TicketGeneratorHandler, error) {
	wire.Build(
		services.NewTicketGenerator,
		NewTicketGeneratorHandler,
	)
	return &TicketGeneratorHandler{}, nil
}
