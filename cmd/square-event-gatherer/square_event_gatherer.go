package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/dv_aws"
	"github.com/timhugh/digitalvenue/dv_aws/dv_dynamodb"
	"github.com/timhugh/digitalvenue/square"
)

type SquareEventGathererHandler struct {
	log      zerolog.Logger
	gatherer square.PaymentGatherer
}

func NewSquareEventGathererHandler(log zerolog.Logger, gatherer square.PaymentGatherer) SquareEventGathererHandler {
	return SquareEventGathererHandler{
		log:      log,
		gatherer: gatherer,
	}
}

func (handler SquareEventGathererHandler) Handle(request events.DynamoDBEvent) (events.DynamoDBEventResponse, error) {
	failures := make([]events.DynamoDBBatchItemFailure, 0)

	for _, record := range request.Records {
		log := handler.log.With().Str("eventID", record.EventID).Logger()

		if record.EventName != "INSERT" {
			log.Warn().Str("eventName", record.EventName).Msg("skipping non-INSERT event")
			continue // Not retryable
		}

		payment, err := buildSquarePayment(record)
		if err != nil {
			log.Error().Err(err).Msg("failed to build square payment")
			continue // Not retryable
		}

		log = handler.log.With().
			Str("squarePaymentID", payment.SquarePaymentID).
			Str("squareMerchantID", payment.SquareMerchantID).
			Str("squareOrderID", payment.SquareOrderID).
			Logger()

		err = handler.gatherer.Gather(payment, log)
		if err != nil {
			log.Error().Err(err).Msg("failed to process payment")

			// TODO: distinguish between retryable and non-retryable errors
			failures = append(failures, events.DynamoDBBatchItemFailure{
				ItemIdentifier: record.EventID,
			})

			continue
		}

		log.Info().Msg("processed payment successfully")
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
