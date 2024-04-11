package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
	"strings"
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
			failures = append(failures, events.DynamoDBBatchItemFailure{
				ItemIdentifier: record.EventID,
			})
			continue
		}

		newImage := record.Change.NewImage
		// TODO: Type could be nil
		if itemType := newImage["Type"].String(); itemType != "SquarePayment" {
			log.Warn().Str("type", itemType).Msg("skipping non-SquarePayment event")
			failures = append(failures, events.DynamoDBBatchItemFailure{
				ItemIdentifier: record.EventID,
			})
			continue
		}

		// TODO: PK and SK could be nil if newImage is empty (like if the stream view type gets changed)
		pk := newImage["PK"].String()
		sk := newImage["SK"].String()
		// TODO: Malformed IDs will panic
		squareMerchantID := strings.Split(pk, "#")[1]
		squarePaymentID := strings.Split(sk, "#")[1]

		log = handler.log.With().
			Str("squarePaymentID", squarePaymentID).
			Logger()

		err := handler.gatherer.Gather(squareMerchantID, squarePaymentID)
		if err != nil {
			log.Error().Err(err).Msg("failed to process payment")

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
