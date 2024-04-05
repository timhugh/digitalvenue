package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/core"
	"github.com/timhugh/digitalvenue/square"
	"github.com/timhugh/digitalvenue/square/webhooks"
)

const squareSignatureHeader = "x-square-hmacsha256-signature"
const squareWebhookNotificationURL = "SQUARE_WEBHOOK_NOTIFICATION_URL"

type handler struct {
	webhookNotificationURL string
	merchantRepo           square.MerchantRepository
	log                    zerolog.Logger
	handlerProvider        webhooks.HandlerProvider
}

func newHandler(merchantRepo square.MerchantRepository, handlerProvider webhooks.HandlerProvider, log zerolog.Logger) handler {
	return handler{
		webhookNotificationURL: core.Getenv(squareWebhookNotificationURL),
		merchantRepo:           merchantRepo,
		log:                    log,
		handlerProvider:        handlerProvider,
	}
}

func (handler handler) handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	webhookEvent, err := webhooks.NewWebhookEvent(request.Body)
	if err != nil {
		handler.log.Warn().Err(err).Msg("Failed to create webhook event")
		return errorResponse("unable to process event")
	}
	log := handler.log.With().
		Str("event_id", webhookEvent.EventID()).
		Str("event", webhookEvent.EventType()).
		Str("merchant_id", webhookEvent.MerchantID()).
		Logger()

	merchant, err := handler.merchantRepo.GetSquareMerchant(webhookEvent.MerchantID())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to find merchant")
		return errorResponse("failed to find merchant with ID '%s'", webhookEvent.MerchantID())
	}

	signature := request.Headers[squareSignatureHeader]
	err = webhooks.Validate(request.Body, handler.webhookNotificationURL, merchant.SquareWebhookSignatureKey, signature)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to validate event")
		return errorResponse("invalid signature: %s", signature)
	}

	eventHandler, err := handler.handlerProvider.GetHandler(webhookEvent.EventType())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get event handler")
		return errorResponse("unknown event type: %s", webhookEvent.EventType())
	}

	err = eventHandler.HandleEvent(webhookEvent)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to handle event")
		return errorResponse("failed to handle event: %s", err.Error())
	}

	log.Info().Msg("Event processed successfully")

	return events.APIGatewayProxyResponse{
		Body:       `{"status": "success"}`,
		StatusCode: 200,
	}, nil
}

func errorResponse(msg string, params ...interface{}) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(`{"error": "%s"}`, fmt.Sprintf(msg, params...)),
		StatusCode: 400,
	}, nil
}
