package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "square-events")
	handler, err := initializeHandler(log)
	if err != nil {
		log.Fatal("Failed to initialize handler")
		os.Exit(1)
	}
	lambda.Start(handler.Handle)
}
