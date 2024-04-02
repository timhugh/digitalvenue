package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
)

type handler struct {
	log      zerolog.Logger
	gatherer square.SquareEventGatherer
}

func newHandler(log zerolog.Logger) handler {
	return handler{log: log}
}

func (handler handler) handle(request events.SQSEvent) (events.SQSEventResponse, error) {
	failures := make([]events.SQSBatchItemFailure, 0)

	for _, record := range request.Records {
		squarePaymentID := record.Body
		err := handler.gatherer.Gather(squarePaymentID)
		if err != nil {
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	response := events.SQSEventResponse{}
	if len(failures) > 0 {
		response.BatchItemFailures = failures
	}
	return response, nil
}
