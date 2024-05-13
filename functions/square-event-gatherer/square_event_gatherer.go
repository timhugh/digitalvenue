package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/timhugh/digitalvenue/util/dv_aws/dv_sqs"
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

func (handler SquareEventGathererHandler) Handle(event events.SQSEvent) (events.SQSEventResponse, error) {
	failures := make([]events.SQSBatchItemFailure, 0)

	for _, record := range event.Records {
		log := handler.log.Sub().AddParam("messageID", record.MessageId)

		var paymentEvent dv_sqs.SquarePaymentCreatedEvent
		err := json.Unmarshal([]byte(record.Body), &paymentEvent)
		if err != nil {
			log.AddParam("error", err.Error()).Error("failed to unmarshal SquarePaymentCreatedEvent")
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		log.AddParams(map[string]interface{}{
			"squarePaymentID":  paymentEvent.SquarePaymentID,
			"squareMerchantID": paymentEvent.SquareMerchantID,
		})

		err = handler.gatherer.Gather(log.NewContext(), paymentEvent.SquareMerchantID, paymentEvent.SquarePaymentID)
		if err != nil {
			log.AddParam("error", err.Error()).Error("failed to process square payment")

			// TODO: distinguish between retryable and non-retryable errors
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})

			continue
		}
	}

	response := events.SQSEventResponse{}
	if len(failures) > 0 {
		response.BatchItemFailures = failures
	}
	return response, nil
}
