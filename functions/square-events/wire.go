//go:build wireinject
// +build wireinject

package main

import (
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/webhooks"
)

func initializeHandler(log *logger.ContextLogger) (*SquareEventsHandler, error) {
	wire.Build(
		dv_aws.DefaultConfig,

		dv_dynamodb.NewClient,
		wire.Bind(new(dv_dynamodb.Client), new(*awsdynamodb.Client)),
		dv_dynamodb.NewRepository,
		wire.Bind(new(square.MerchantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(square.PaymentRepository), new(*dv_dynamodb.Repository)),

		webhooks.NewHandlerProvider,
		webhooks.NewPaymentCreatedHandler,

		NewSquareEventsHandler,
	)
	return &SquareEventsHandler{}, nil
}
