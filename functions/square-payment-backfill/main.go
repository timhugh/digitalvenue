package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"github.com/timhugh/digitalvenue/util/square"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "square-payment-backfill")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	handler, err := initializeHandler(log)
	if err != nil {
		log.AddParam("error", err.Error()).Fatal("Failed to initialize handler")
		os.Exit(1)
	}
	lambda.Start(handler.Handle)
}

type SquarePaymentBackfillRequest struct {
	BackfillDate string `json:"backfill_date"`
	TenantID     string `json:"tenant_id"`
}

type SquarePaymentBackfillHandler struct {
	log       logger.ContextLogger
	squareAPI square.APIClient
}

func NewSquarePaymentBackfillHandler(
	squareAPI square.APIClient,
) (*SquarePaymentBackfillHandler, error) {
	return &SquarePaymentBackfillHandler{
		squareAPI: squareAPI,
	}, nil
}

func (h *SquarePaymentBackfillHandler) Handle(apiRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var request SquarePaymentBackfillRequest
	err := json.Unmarshal([]byte(apiRequest.Body), &request)
	if err != nil {
		h.log.AddParam("error", err.Error()).Error("Error unmarshalling request body")
	}

	return events.APIGatewayProxyResponse{
		Body:       `{"status": "success"}`,
		StatusCode: 200,
	}, nil
}
