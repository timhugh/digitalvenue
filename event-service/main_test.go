package main

import (
	"fmt"
	"github.com/matryer/is"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var webhookEventRawJson, _ = os.ReadFile("test-event.json")
var webhookEventJson = string(webhookEventRawJson)

func TestHandler(t *testing.T) {
	t.Skip("implementation WIP")
	is := is.New(t)

	testCases := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		expectedStatus int
		expectedBody   string
		expectedError  error
	}{
		{
			name: "basic success",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJson,
				Headers: map[string]string{
					squareSignatureHeader: "BJiV0YFeQNzdtt+z/J9xsk3omhM4wLOhTfK76ZL5thc=",
					"Content-Type":        "application/json",
				},
			},
			expectedStatus: 200,
			expectedBody:   "",
			expectedError:  nil,
		},
		{
			name: "bad json",
			request: events.APIGatewayProxyRequest{
				Body: "this isn't even json",
			},
			expectedStatus: 400,
			expectedBody:   "",
			expectedError:  fmt.Errorf("invalid request"),
		},
		{
			name: "incorrect signature",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJson,
				Headers: map[string]string{
					squareSignatureHeader: "not the right signature",
				},
			},
			expectedStatus: 400,
			expectedBody:   "",
			expectedError:  fmt.Errorf("invalid signature"),
		},
		{
			name: "unknown event type",
			request: events.APIGatewayProxyRequest{
				Body: `{"type": "not.a.real.event"}`,
			},
			expectedStatus: 400,
			expectedBody:   "",
			expectedError:  fmt.Errorf("unknown event type"),
		},
	}

	config := EventServiceConfig{
		WebhookUrl: "http://localhost:8080/events",
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//mockEventHandler := mock.Mock[core.EventHandler]()
			handler := handler(config)

			response, err := handler(testCase.request)

			is.Equal(err, testCase.expectedError)
			is.Equal(response.StatusCode, testCase.expectedStatus)
			is.Equal(response.Body, testCase.expectedBody)
		})
	}
}
