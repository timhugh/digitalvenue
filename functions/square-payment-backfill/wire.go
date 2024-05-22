//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square/squareapi"
)

func initializeHandler(logger *logger.ContextLogger) (*SquarePaymentBackfillHandler, error) {
	wire.Build(
		squareapi.NewClient,

		NewSquarePaymentBackfillHandler,
	)
	return &SquarePaymentBackfillHandler{}, nil
}
