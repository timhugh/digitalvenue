package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/persistence"
	"github.com/timhugh/digitalvenue/persistence/dynamodb"
	"github.com/timhugh/digitalvenue/square/webhooks"
	"os"
)

const squareSignatureHeader = "x-square-hmacsha256-signature"

type EventServiceConfig struct {
	WebhookUrl string
}

func handler(config EventServiceConfig, merchantRepo persistence.MerchantRepo) func(ctx events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		webhookEvent, err := webhooks.NewWebhookEvent(request.Body)
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Err(err).
				Msg("Failed to create webhook event")

			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "Unable to process event: %s"}`, err.Error()),
			}, nil
		}

		merchant, err := merchantRepo.FindMerchantBySquareMerchantId(webhookEvent.MerchantId())
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Str("event_id", webhookEvent.EventId()).
				Str("event", webhookEvent.EventType()).
				Str("merchant_id", webhookEvent.MerchantId()).
				Err(err).
				Msg("Failed to find merchant")

			// TODO: unknown merchant isn't entirely accurate
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "unknown merchant: %s"}`, webhookEvent.MerchantId()),
			}, nil
		}

		signature := request.Headers[squareSignatureHeader]
		err = webhooks.Validate(request.Body, config.WebhookUrl, merchant.SquareWebhookSignatureKey, signature)
		if err != nil {
			log.Warn().
				Str("service", "events-service").
				Str("event_id", webhookEvent.EventId()).
				Str("event", webhookEvent.EventType()).
				Str("merchant_id", webhookEvent.MerchantId()).
				Err(err).
				Msg("Failed to validate event")

			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "invalid signature: %s"}`, signature),
			}, nil
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

			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf(`{"error": "unknown event type: %s"}`, webhookEvent.EventType()),
			}, nil
		}

		eventHandler.HandleEvent(webhookEvent)

		return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 200,
		}, nil
	}
}

func main() {
	log.Info().Msg("Starting events-service")

	config := EventServiceConfig{
		WebhookUrl: os.Getenv("WEBHOOK_NOTIFICATION_URL"),
	}
	merchantRepo, err := dynamodb.NewMerchantRepo(os.Getenv("MERCHANTS_TABLE"))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create merchant repo")
	}

	handler := handler(config, merchantRepo)
	lambda.Start(handler)
}
