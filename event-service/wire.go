//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/db/dynamodb"
)

func initializeHandler() (handler, error) {
	wire.Build(
		newHandler,
		newEventServiceConfig,
		dynamodb.NewMerchantsRespository,
		dynamodb.NewMerchantsRepositoryConfig,
	)
	return handler{}, nil
}
