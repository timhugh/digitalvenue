package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/dv_aws"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
)

type SquareEventGathererHandler struct {
	log      *logger.ContextLogger
	gatherer square.PaymentGatherer
}

func NewSquareEventGathererHandler(log *logger.ContextLogger, gatherer square.PaymentGatherer) SquareEventGathererHandler {
	return SquareEventGathererHandler{
		log:      log,
		gatherer: gatherer,
	}
}

func (handler SquareEventGathererHandler) Handle(request events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	failures := make([]events.DynamoDBBatchItemFailure, 0)

	for _, record := range request.Records {
		log := handler.log.Sub().AddParam("eventID", record.EventID)

		if record.EventName != "INSERT" {
			log.Warn("skipping non-INSERT event '%s'", record.EventName)
			continue // Not retryable
		}

		payment, err := buildSquarePayment(record)
		if err != nil {
			log.Error("failed to build square payment: %s", err)
			continue // Not retryable
		}

		log.AddParams(map[string]interface{}{
			"squarePaymentID":  payment.SquarePaymentID,
			"squareMerchantID": payment.SquareMerchantID,
			"squareOrderID":    payment.SquareOrderID,
		})

		err = handler.gatherer.Gather(log.NewContext(), payment)
		if err != nil {
			log.Error("Failed to process payment: %s", err)

			// TODO: distinguish between retryable and non-retryable errors
			failures = append(failures, events.DynamoDBBatchItemFailure{
				ItemIdentifier: record.EventID,
			})

			continue
		}

		log.Info("Processed payment successfully")
	}

	response := events.DynamoDBEventResponse{}
	if len(failures) > 0 {
		response.BatchItemFailures = failures
	}
	return response, nil
}

func buildSquarePayment(record events.DynamoDBEventRecord) (*square.Payment, error) {
	attrs, err := dv_aws.GetImageAttributes("SquarePayment", record.Change.NewImage, "PK", "SK", "SquareOrderID")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get square payment attributes from dynamodb event new image")
	}

	squarePaymentID, err := dv_dynamodb.UnprefixID(attrs["SK"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to unprefix square payment ID")
	}
	squareMerchantID, err := dv_dynamodb.UnprefixID(attrs["PK"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to unprefix square merchant ID")
	}

	return &square.Payment{
		SquarePaymentID:  squarePaymentID,
		SquareMerchantID: squareMerchantID,
		SquareOrderID:    attrs["SquareOrderID"],
	}, nil
}
