package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/core/services"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
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
}

func NewTicketGeneratorHandler(log *logger.ContextLogger, generator *services.TicketGenerator) *TicketGeneratorHandler {
	return &TicketGeneratorHandler{
		log:       log,
		generator: generator,
	}
}

func (handler *TicketGeneratorHandler) Handle(request events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	failures := make([]events.DynamoDBBatchItemFailure, 0)

	handler.log.Info("Processing batch of records")

	// TODO: process records concurrently
	for _, record := range request.Records {
		log := handler.log.Sub().AddParams(map[string]interface{}{
			"eventID":   record.EventID,
			"eventName": record.EventName,
		})
		log.Info("Processing record")

		order, err := buildOrderFromEvent(record)
		if err != nil {
			log.Error("Failed to build order from event: %s", err)
			continue // not retryable
		}

		log.AddParams(map[string]interface{}{
			"orderID":  order.ID,
			"tenantID": order.TenantID,
		})

		err = handler.generator.GenerateTickets(log.NewContext(), order)
		if err != nil {
			log.Error("Failed to generate tickets: %s", err)
			failures = append(failures, events.DynamoDBBatchItemFailure{ItemIdentifier: record.EventID})
			continue // not retryable
		}

		log.Debug("Successfully processed record")
	}

	if len(failures) > 0 {
		handler.log.Error("Failed to process some records")
		return events.DynamoDBEventResponse{BatchItemFailures: failures}, nil
	}

	handler.log.Debug("Successfully processed all records")
	return events.DynamoDBEventResponse{}, nil
}

func buildOrderFromEvent(record events.DynamoDBEventRecord) (*core.Order, error) {
	newImage := record.Change.NewImage
	if newImage == nil {
		return nil, errors.New("record does not contain NewImage")
	}

	imageType, ok := newImage["Type"]
	if !ok {
		return nil, errors.New("record does not contain Type")
	}
	if imageType.String() != "Order" {
		return nil, errors.New("record is not of type Order")
	}

	var order core.Order

	pk, ok := newImage["PK"]
	if !ok {
		return nil, errors.New("record does not contain PK")
	}
	tenantID, err := dv_dynamodb.UnprefixID(pk.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to unprefix TenantID")
	}
	order.TenantID = tenantID

	sk, ok := newImage["SK"]
	if !ok {
		return nil, errors.New("record does not contain SK")
	}
	orderID, err := dv_dynamodb.UnprefixID(sk.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to unprefix OrderID")
	}
	order.ID = orderID

	customerID, ok := newImage["CustomerID"]
	if !ok {
		return nil, errors.New("record does not contain CustomerID")
	}
	order.CustomerID = customerID.String()

	metaRaw, ok := newImage["Meta"]
	if ok {
		metaMap := metaRaw.Map()
		order.Meta = make(map[string]string, len(metaMap))
		for k, v := range metaMap {
			order.Meta[k] = v.String()
		}
	} else {
		order.Meta = make(map[string]string)
	}

	orderItemsRaw, ok := newImage["OrderItems"]
	if !ok {
		return nil, errors.New("record does not contain OrderItems")
	}
	orderItemsRawList := orderItemsRaw.List()
	order.Items = make([]core.OrderItem, len(orderItemsRawList))
	for i, itemRaw := range orderItemsRawList {
		var item core.OrderItem

		itemMap := itemRaw.Map()

		itemID, ok := itemMap["ItemID"]
		if !ok {
			return nil, errors.New("OrderItem does not contain ItemID")
		}
		item.ID = itemID.String()

		name, ok := itemMap["Name"]
		if !ok {
			return nil, errors.New("OrderItem does not contain Name")
		}
		item.Name = name.String()

		metaRaw, ok := itemMap["Meta"]
		if ok {
			metaMap := metaRaw.Map()
			item.Meta = make(map[string]string, len(metaMap))
			for k, v := range metaMap {
				item.Meta[k] = v.String()
			}
		} else {
			item.Meta = make(map[string]string)
		}

		order.Items[i] = item
	}

	return &order, nil
}
