package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "ticket-generator")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	handler := initializeHandler(log)
	lambda.Start(handler.Handle)
}

type TicketMailerHandler struct {
	log *logger.ContextLogger
}

func NewTicketMailerHandler(log *logger.ContextLogger) *TicketMailerHandler {
	return &TicketMailerHandler{
		log: log,
	}
}

func (h *TicketMailerHandler) Handle(event interface{}) {
	h.log.Info("Received event: %#v", event)
}
