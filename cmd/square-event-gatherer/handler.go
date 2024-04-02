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

func (handler handler) handle(request events.SQSEvent) (events.SQSEventResponse, error) {

	fmt.Printf("Received request: %v\n", request)

	return events.SQSEventResponse{}, nil
}
