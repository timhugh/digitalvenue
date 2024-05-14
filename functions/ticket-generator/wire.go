//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/core/services"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_s3"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_sqs"
	"github.com/timhugh/digitalvenue/util/logger"
)

func initializeHandler(logger *logger.ContextLogger) (*TicketGeneratorHandler, error) {
	wire.Build(
		services.NewTicketGenerator,
		NewTicketGeneratorHandler,

		dv_aws.DefaultConfig,

		dv_s3.NewClient,
		dv_s3.NewS3QRStorage,
		wire.Bind(new(core.QRCodeStore), new(*dv_s3.S3QRStorage)),

		dv_dynamodb.NewClient,
		dv_dynamodb.NewRepository,
		wire.Bind(new(core.TicketRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dv_dynamodb.Repository)),

		dv_sqs.NewClient,
		dv_sqs.NewOrderProcessedQueue,
		wire.Bind(new(core.OrderProcessedQueue), new(*dv_sqs.OrderProcessedQueue)),
	)
	return &TicketGeneratorHandler{}, nil
}
