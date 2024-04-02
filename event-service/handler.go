package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/timhugh/digitalvenue/db"
	webhooks2 "github.com/timhugh/digitalvenue/services/webhooks"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

const squareSignatureHeader = "x-square-hmacsha256-signature"

type handler struct {
	config          eventServiceConfig
	merchantRepo    db.MerchantsRepository
	log             zerolog.Logger
	handlerProvider webhooks2.HandlerProvider
}

func newHandler(config eventServiceConfig, merchantRepo db.MerchantsRepository, handlerProvider webhooks2.HandlerProvider) handler {
	return handler{
		config:          config,
		merchantRepo:    merchantRepo,
		log:             log.With().Str("service", "event-service").Logger(),
		handlerProvider: handlerProvider,
	}
}

func (h handler) handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	webhookEvent, err := webhooks.NewWebhookEvent(request.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create webhook event")

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Unable to process event: %s"}`, err.Error()),
		}, nil
	}
	log := log.With().
		Str("event_id", webhookEvent.EventId()).
		Str("event", webhookEvent.EventType()).
		Str("merchant_id", webhookEvent.MerchantId()).
		Logger()

	merchant, err := h.merchantRepo.FindMerchantBySquareMerchantId(webhookEvent.MerchantId())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to find merchant")

		// TODO: unknown merchant isn't entirely accurate
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "unknown merchant: %s"}`, webhookEvent.MerchantId()),
		}, nil
	}

	signature := request.Headers[squareSignatureHeader]
	err = webhooks.Validate(request.Body, h.config.webhookUrl, merchant.SquareWebhookSignatureKey, signature)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to validate event")

		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "invalid signature: %s"}`, signature),
		}, nil
	}

	eventHandler, err := h.handlerProvider.GetHandler(webhookEvent.EventType())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get event handler")

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
