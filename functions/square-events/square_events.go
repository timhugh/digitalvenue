package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"github.com/timhugh/digitalvenue/util/square/webhooks"
)

const squareSignatureHeader = "x-square-hmacsha256-signature"
const squareWebhookNotificationURL = "SQUARE_WEBHOOK_NOTIFICATION_URL"

type SquareEventsHandler struct {
	webhookNotificationURL string
	merchantRepo           square.MerchantRepository
	log                    *logger.ContextLogger
	handlerProvider        webhooks.HandlerProvider
}

func NewSquareEventsHandler(merchantRepo square.MerchantRepository, handlerProvider webhooks.HandlerProvider, log *logger.ContextLogger) (*SquareEventsHandler, error) {
	squareWebhookNotificationURL, err := core.RequireEnv(squareWebhookNotificationURL)
	if err != nil {
		return nil, err
	}

	return &SquareEventsHandler{
		webhookNotificationURL: squareWebhookNotificationURL,
		merchantRepo:           merchantRepo,
		log:                    log,
		handlerProvider:        handlerProvider,
	}, nil
}

func (handler *SquareEventsHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log := handler.log.Sub().AddParams(map[string]interface{}{
		"requestID": request.RequestContext.RequestID,
		"apiStage":  request.RequestContext.Stage,
	})

	webhookEvent, err := webhooks.NewWebhookEvent(request.Body)
	if err != nil {
		log.AddParam("error", err.Error()).Error("Failed to create webhook event")
		return errorResponse("unable to process event: %s", err.Error())
	}

	log.AddParams(map[string]interface{}{
		"event_id":    webhookEvent.EventID(),
		"event":       webhookEvent.EventType(),
		"merchant_id": webhookEvent.MerchantID(),
	})
	log.Info("Handling webhook event")

	merchant, err := handler.merchantRepo.GetSquareMerchant(webhookEvent.MerchantID())
	if err != nil {
		log.AddParam("error", err.Error()).Error("Failed to find merchant")
		return errorResponse("failed to find merchant with ID '%s'", webhookEvent.MerchantID())
	}

	webhookEvent.SetTenantID(merchant.TenantID)
	log.AddParam("tenant_id", merchant.TenantID)

	signature := request.Headers[squareSignatureHeader]
	err = webhooks.Validate(request.Body, handler.webhookNotificationURL, merchant.SquareWebhookSignatureKey, signature)
	if err != nil {
		log.AddParam("error", err.Error()).Error("Failed to validate event")
		return errorResponse("invalid signature: %s", signature)
	}

	eventHandler, err := handler.handlerProvider.GetHandler(webhookEvent.EventType())
	if err != nil {
		log.AddParam("error", err.Error()).Error("Failed to get event handler")
		return errorResponse("unknown event type: %s", webhookEvent.EventType())
	}

	ctx := logger.Attach(context.Background(), log)
	err = eventHandler.HandleEvent(ctx, webhookEvent)
	if err != nil {
		log.AddParam("error", err.Error()).Error("Failed to handle event")
		return errorResponse("failed to handle event: %s", err.Error())
	}

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
