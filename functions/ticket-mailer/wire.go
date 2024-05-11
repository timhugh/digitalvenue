//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_s3"
	"github.com/timhugh/digitalvenue/util/logger"
)

func initializeHandler(log *logger.ContextLogger) (*TicketMailerHandler, error) {
	wire.Build(
		NewTicketMailerHandler,

		dv_aws.DefaultConfig,
		dv_dynamodb.NewClient,
		dv_dynamodb.NewRepository,
		wire.Bind(new(core.TenantRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.OrderRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.CustomerRepository), new(*dv_dynamodb.Repository)),
		wire.Bind(new(core.TicketRepository), new(*dv_dynamodb.Repository)),

		dv_s3.NewClient,
		dv_s3.NewS3TemplateStore,
		wire.Bind(new(core.TemplateStore), new(*dv_s3.TemplateStore)),
	)
	return &TicketMailerHandler{}, nil
}
