package handlers

import (
	"errors"
	"github.com/matryer/is"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/rs/zerolog"
	"github.com/timhugh/digitalvenue/square"
	squarewebhooks "github.com/timhugh/digitalvenue/square/webhooks"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var webhookEventRawJSON, _ = os.ReadFile("square_event_test_body.json")
var webhookEventJSON = string(webhookEventRawJSON)

const goodSignature = "PlQk/ad1RoIYA2at/LC21DGJzKz0J/xAZ8KniOf4ouo="

func TestSquareEventsHandler(t *testing.T) {
	testCases := []struct {
		name string
		// given
		request            events.APIGatewayProxyRequest
		merchant           square.Merchant
		merchantFetchError error
		// then
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "basic success",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJSON,
				Headers: map[string]string{
					squareSignatureHeader: goodSignature,
					"Content-Type":        "application/json",
				},
			},
			merchant:       square.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 200,
			expectedBody:   `{"status": "success"}`,
		},
		{
			name: "bad json",
			request: events.APIGatewayProxyRequest{
				Body: "this isn't even json",
			},
			expectedStatus: 400,
			expectedBody:   `{"error": "unable to process event: failed to unmarshal webhook event metadata"}`,
		},
		{
			name: "unknown merchant",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJSON,
				Headers: map[string]string{
					squareSignatureHeader: goodSignature,
				},
			},
			merchantFetchError: errors.New("who dat"),
			expectedStatus:     400,
			expectedBody:       `{"error": "failed to find merchant with ID 'merchant_id'"}`,
		},
		{
			name: "incorrect signature",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJSON,
				Headers: map[string]string{
					squareSignatureHeader: "invalid_signature",
				},
			},
			merchant:       square.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 400,
			expectedBody:   `{"error": "invalid signature: invalid_signature"}`,
		},
		{
			name: "unknown event type",
			request: events.APIGatewayProxyRequest{
				Body: `{"type": "not.a.real.event"}`,
			},
			merchant:       square.Merchant{SquareWebhookSignatureKey: "signature_key"},
			expectedStatus: 400,
			expectedBody:   `{"error": "unable to process event: unknown event type: not.a.real.event"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)
			mock.SetUp(t)

			err := os.Setenv("SQUARE_WEBHOOK_NOTIFICATION_URL", "https://example.com")
			is.NoErr(err)

			mockMerchantRepo := mock.Mock[square.MerchantRepository]()
			mock.WhenDouble(mockMerchantRepo.GetSquareMerchant(mock.Any[string]())).ThenReturn(testCase.merchant, testCase.merchantFetchError)

			mockHandlerProvider := mock.Mock[squarewebhooks.HandlerProvider]()
			mockHandler := mock.Mock[squarewebhooks.EventHandler]()
			mock.WhenDouble(mockHandlerProvider.GetHandler(mock.Any[string]())).ThenReturn(mockHandler, nil)
			mock.WhenSingle(mockHandler.HandleEvent(mock.Any[squarewebhooks.WebhookEvent[any]]())).ThenReturn(nil)

			handler, err := NewSquareEventsHandler(mockMerchantRepo, mockHandlerProvider, zerolog.Logger{})
			is.NoErr(err)

			response, err := handler.Handle(testCase.request)
			is.NoErr(err)
			is.Equal(response.StatusCode, testCase.expectedStatus)
			is.Equal(response.Body, testCase.expectedBody)
		})
	}
}
