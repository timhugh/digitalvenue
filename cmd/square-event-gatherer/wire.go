//go:build wireinject
// +build wireinject

package main

import (
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/dv_aws"
	"github.com/timhugh/digitalvenue/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/logger"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/squareapi"
)

func initializeHandler(log *logger.ContextLogger) (SquareEventGathererHandler, error) {
	wire.Build(
		square.NewPaymentGatherer,

		dv_aws.DefaultConfig,

		dv_dynamodb.NewClient,
		wire.Bind(new(dv_dynamodb.Client), new(*awsdynamodb.Client)),
		dv_dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dv_dynamodb.Repository)),

		squareapi.NewClient,
		wire.Bind(new(square.APIClient), new(*squareapi.Client)),

		NewSquareEventGathererHandler,
	)
	return SquareEventGathererHandler{}, nil
}
