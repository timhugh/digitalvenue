package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/core/services"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_sqs"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "ticket-generator")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	handler, err := initializeHandler(log)
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to initialize handler")
		os.Exit(1)
	}
	lambda.Start(handler.Handle)
}

type TicketGeneratorHandler struct {
	log       *logger.ContextLogger
	generator *services.TicketGenerator
	orderRepo core.OrderRepository
}

func NewTicketGeneratorHandler(
	log *logger.ContextLogger,
	generator *services.TicketGenerator,
	orderRepo core.OrderRepository,
) *TicketGeneratorHandler {
	return &TicketGeneratorHandler{
		log:       log,
		generator: generator,
		orderRepo: orderRepo,
	}
}

func (handler *TicketGeneratorHandler) Handle(event events.SQSEvent) (events.SQSEventResponse, error) {
	failures := make([]events.SQSBatchItemFailure, 0)

	handler.log.Info("Processing event batch")

	// TODO: process records concurrently
	for _, record := range event.Records {
		log := handler.log.Sub().AddParams(map[string]interface{}{
			"messageID": record.MessageId,
		})
		log.Info("Processing record")

		var orderEvent dv_sqs.OrderCreatedEvent
		err := json.Unmarshal([]byte(record.Body), &orderEvent)
		if err != nil {
			log.AddParam("error", err.Error()).Fatal("Failed to unmarshal OrderCreatedEvent")
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		order, err := handler.orderRepo.GetOrder(orderEvent.TenantID, orderEvent.OrderID)
		if err != nil {
			log.AddParam("error", err.Error()).Fatal("Failed to get order")
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		log.AddParams(map[string]interface{}{
			"orderID":  order.ID,
			"tenantID": order.TenantID,
		})

		err = handler.generator.GenerateTickets(log.NewContext(), order)
		if err != nil {
			log.Error("Failed to generate tickets: %s", err)
			failures = append(failures, events.SQSBatchItemFailure{ItemIdentifier: record.MessageId})
			continue
		}

		log.Debug("Successfully processed record")
	}

	if len(failures) > 0 {
		handler.log.Error("Failed to process some records")
		return events.SQSEventResponse{BatchItemFailures: failures}, nil
	}

	handler.log.Debug("Successfully processed all records")
	return events.SQSEventResponse{}, nil
}
