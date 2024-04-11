package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/dv_aws/dv_dynamodb"
)

func main() {
	logger := zerolog.Logger{}.With().Str("service", "ticket-generator").Logger()
	handler, err := initializeHandler(logger)
	if err != nil {
		logger.Fatal().Err(err).Str("service", "ticket-generator").Msg("Failed to initialize handler")
	}
	lambda.Start(handler.Handle)
}

type TicketGeneratorHandler struct {
	logger zerolog.Logger
}

func NewTicketGeneratorHandler(logger zerolog.Logger) *TicketGeneratorHandler {
	return &TicketGeneratorHandler{
		logger: logger,
	}
}

func (handler *TicketGeneratorHandler) Handle(request events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	//ctx := context.Background()
	//failures := make([]events.DynamoDBBatchItemFailure, 0)

	for _, record := range request.Records {
		logger := handler.logger.With().Str("eventID", record.EventID).Str("eventName", record.EventName).Logger()
		logger.Info().Msg("Processing record")

		order, err := buildOrderFromEvent(record)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to build order from event")
			continue // not retryable
		}

		logger = logger.With().
			Str("orderID", order.ID).
			Str("tenantID", order.TenantID).
			Logger()

		logger.Info().Msg("Retrieved order from event")
	}

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
