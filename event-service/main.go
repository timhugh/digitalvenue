package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

const squareSignatureHeader = "x-square-hmacsha256-signature"

type EventServiceConfig struct {
	WebhookUrl string
}

func handler(config EventServiceConfig) func(ctx events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		webhookEvent, err := webhooks.NewWebhookEvent(request.Body)
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Err(err).
				Msg("Failed to unmarshal request body")

			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("malformed request json")
		}

		// TODO: find merchant to get actual signature key

		err = webhooks.Validate(request.Body, config.WebhookUrl, "signature_key", request.Headers[squareSignatureHeader])
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Str("event_id", webhookEvent.EventId()).
				Str("event", webhookEvent.EventType()).
				Str("merchant_id", webhookEvent.MerchantId()).
				Err(err).
				Msg("Failed to validate event")

			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("invalid signature")
		}

		eventHandler, err := core.GetHandler(webhookEvent.EventType())
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Str("event_id", webhookEvent.EventId()).
				Str("event", webhookEvent.EventType()).
				Str("merchant_id", webhookEvent.MerchantId()).
				Err(err).
				Msg("Failed to get event handler")

			return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("unknown event type")
		}

		eventHandler.HandleEvent(webhookEvent)

		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 200,
		}, nil
	}
}

func main() {
	config := EventServiceConfig{
		// TODO
		WebhookUrl: "http://localhost:8080/events",
	}

	handler := handler(config)
	lambda.Start(handler)
}
