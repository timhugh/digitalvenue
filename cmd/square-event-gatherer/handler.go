package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

type handler struct {
	log zerolog.Logger
}

func newHandler(log zerolog.Logger) handler {
	return handler{log: log}
}

func (handler handler) handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("Received request: %v\n", request)

	return events.APIGatewayProxyResponse{}, nil
}
