//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/logger"
)

func initializeHandler(log *logger.ContextLogger) *TicketMailerHandler {
	wire.Build(
		NewTicketMailerHandler,
	)
	return &TicketMailerHandler{}
}
