package main

import (
	"github.com/matryer/is"
	"github.com/rs/zerolog"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

var webhookEventRawJSON, _ = os.ReadFile("test-event.json")
var webhookEventJSON = string(webhookEventRawJSON)

const goodSignature = "/p9MrQ6sTzL2iuGBPa5YoadntDIMv5ms+ihDe3MLoLc="

func TestHandler(t *testing.T) {
	testCases := []struct {
		name string
		// given
		request events.APIGatewayProxyRequest
		// then
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "basic success",
			request: events.APIGatewayProxyRequest{
				Body: webhookEventJSON,
			},
			expectedStatus: 200,
			expectedBody:   `{"status": "success"}`,
		},
	}

	log := zerolog.Logger{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			is := is.New(t)

			handler := newHandler(log)

			response, err := handler.handle(testCase.request)
			is.NoErr(err)
			is.Equal(response.StatusCode, testCase.expectedStatus)
			is.Equal(response.Body, testCase.expectedBody)
		})
	}
}
