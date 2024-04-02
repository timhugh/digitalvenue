package main

import (
	"fmt"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square/db"
	squarewebhooks "github.com/timhugh/digitalvenue/square/webhooks"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var webhookEventRawJson, _ = os.ReadFile("test-event.json")
var webhookEventJson = string(webhookEventRawJson)

const goodSignature = "/p9MrQ6sTzL2iuGBPa5YoadntDIMv5ms+ihDe3MLoLc="

func TestHandler(t *testing.T) {
	testCases := []struct {
		name string
		// given
		request            events.APIGatewayProxyRequest
		merchant           db.Merchant
		merchantFetchError error
		// then
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "basic success",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJson,
				Headers: map[string]string{
					squareSignatureHeader: goodSignature,
					"Content-Type":        "application/json",
				},
			},
			merchant:       db.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 200,
			expectedBody:   `{"status": "success"}`,
		},
		{
			name: "bad json",
			request: events.APIGatewayProxyRequest{
				Body: "this isn't even json",
			},
			expectedStatus: 400,
			expectedBody:   `{"error": "unable to process event: malformed request json"}`,
		},
		{
			name: "unknown merchant",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJson,
				Headers: map[string]string{
					squareSignatureHeader: goodSignature,
				},
			},
			merchantFetchError: fmt.Errorf("who dat"),
			expectedStatus:     400,
			expectedBody:       `{"error": "failed to find merchant with ID 'merchant_id'"}`,
		},
		{
			name: "incorrect signature",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJson,
				Headers: map[string]string{
					squareSignatureHeader: "not the right signature",
				},
			},
			merchant:       db.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 400,
			expectedBody:   `{"error": "invalid signature: not the right signature"}`,
		},
		{
			name: "unknown event type",
			request: events.APIGatewayProxyRequest{
				Body: `{"type": "not.a.real.event"}`,
			},
			merchant:       db.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 400,
			expectedBody:   `{"error": "unable to process event: unknown event type: not.a.real.event"}`,
		},
	}

	config := eventServiceConfig{
		webhookUrl: "http://localhost:8080/events",
	}

	log := zerolog.Logger{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)
			mock.SetUp(t)

			mockMerchantRepo := mock.Mock[db.MerchantsRepository]()
			mockHandlerProvider := mock.Mock[squarewebhooks.HandlerProvider]()
			mockHandler := mock.Mock[squarewebhooks.EventHandler]()
			mock.WhenDouble(mockMerchantRepo.FindMerchantBySquareMerchantId(mock.Any[string]())).ThenReturn(testCase.merchant, testCase.merchantFetchError)
			mock.WhenDouble(mockHandlerProvider.GetHandler(mock.Any[string]())).ThenReturn(mockHandler, nil)
			mock.WhenSingle(mockHandler.HandleEvent(mock.Any[squarewebhooks.WebhookEvent[any]]())).ThenReturn(nil)
			handler := newHandler(config, mockMerchantRepo, mockHandlerProvider, log)

			response, err := handler.handle(testCase.request)
			is.NoErr(err)
			is.Equal(response.StatusCode, testCase.expectedStatus)
			is.Equal(response.Body, testCase.expectedBody)
		})
	}
}
