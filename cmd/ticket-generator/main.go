package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	handler, err := initializeHandler()
	if err != nil {
		log.Fatal().Err(err).Str("service", "ticket-generator").Msg("Failed to initialize handler")
	}
	lambda.Start(handler.Handle)
}

type TicketGeneratorHandler struct {
	log zerolog.Logger
}

func NewTicketGeneratorHandler(log zerolog.Logger) *TicketGeneratorHandler {
	return &TicketGeneratorHandler{
		log: log,
	}
}

func (handler *TicketGeneratorHandler) Handle() {

}
