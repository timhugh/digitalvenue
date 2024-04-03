package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
)

type handler struct {
	log      zerolog.Logger
	gatherer square.EventGatherer
}

func newHandler(log zerolog.Logger, gatherer square.EventGatherer) handler {
	return handler{
		log:      log,
		gatherer: gatherer,
	}
}

func (handler handler) handle(request events.SQSEvent) (events.SQSEventResponse, error) {
	failures := make([]events.SQSBatchItemFailure, 0)

	for _, record := range request.Records {
		squarePaymentID := record.Body
		log := handler.log.With().
			Str("messageId", record.MessageId).
			Str("squarePaymentID", squarePaymentID).
			Logger()

		err := handler.gatherer.Gather(squarePaymentID)
		if err != nil {
			log.Error().Err(err).Msg("failed to process payment")

			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})

			continue
		}

		log.Info().Msg("processed payment successfully")
	}

	response := events.SQSEventResponse{}
	if len(failures) > 0 {
		response.BatchItemFailures = failures
	}
	return response, nil
}
