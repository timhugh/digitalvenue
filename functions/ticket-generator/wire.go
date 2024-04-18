//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/core/services"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_s3"
	"github.com/timhugh/digitalvenue/util/logger"
)

func initializeHandler(logger *logger.ContextLogger) (*TicketGeneratorHandler, error) {
	wire.Build(
		services.NewTicketGenerator,
		NewTicketGeneratorHandler,

		dv_aws.DefaultConfig,

		dv_s3.NewClient,
		wire.Bind(new(dv_s3.Client), new(*s3.Client)),
		dv_s3.NewS3QRStorage,
		wire.Bind(new(core.QRCodeStorer), new(*dv_s3.S3QRStorage)),

		dv_dynamodb.NewClient,
		wire.Bind(new(dv_dynamodb.Client), new(*dynamodb.Client)),
		dv_dynamodb.NewRepository,
		wire.Bind(new(core.TicketRepository), new(*dv_dynamodb.Repository)),
	)
	return &TicketGeneratorHandler{}, nil
}
