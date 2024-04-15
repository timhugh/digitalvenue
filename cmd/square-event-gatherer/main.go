package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/logger"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "square-event-gatherer")
	handler, err := initializeHandler(log)
	if err != nil {
		log.Fatal("Failed to initialize handler")
		os.Exit(1)
	}
	lambda.Start(handler.Handle)
}
