//go:build wireinject
// +build wireinject

package main

import (
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/db/dynamodb"
	"github.com/timhugh/digitalvenue/services/webhooks"
)

func NewPaymentsRepository(config dynamodb.PaymentsRepositoryConfig, client *awsdynamodb.Client) db.PaymentsRepository {
	return dynamodb.NewPaymentsRepository(config, client)
}

func NewMerchantsRepository(config dynamodb.MerchantsRepositoryConfig, client *awsdynamodb.Client) db.MerchantsRepository {
	return dynamodb.NewMerchantsRepository(config, client)
}

func initializeHandler() (handler, error) {
	wire.Build(
		newHandler,
		newEventServiceConfig,
		dynamodb.NewConfig,
		dynamodb.NewClient,
		dynamodb.NewMerchantsRepositoryConfig,
		NewMerchantsRepository,
		dynamodb.NewPaymentsRepositoryConfig,
		NewPaymentsRepository,
		webhooks.NewPaymentCreatedService,
		webhooks.NewHandlerProvider,
	)
	return handler{}, nil
}
