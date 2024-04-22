package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/logger"
	"log/syslog"
	"os"
)

func main() {
	log := logger.Default().AddParam("service", "papertrail-log-streamer")
	env, err := core.RequireEnv("ENVIRONMENT")
	if err != nil {
		log.AddParam("error", err).Fatal("Failed to determine application environment")
		os.Exit(1)
	}
	log.AddParam("environment", env)

	paperTrailPushURL, err := core.RequireEnv("PAPERTRAIL_PUSH_URL")
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	w, err := syslog.Dial("udp", paperTrailPushURL, syslog.LOG_EMERG|syslog.LOG_KERN, "digitalvenue")
	if err != nil {
		log.Fatal("Failed to connect to syslog: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := w.Close()
		if err != nil {
			log.Warn("Error while closing log stream: %s", err.Error())
		}
	}()

	lambda.Start(func(ctx context.Context, event events.CloudwatchLogsEvent) {
		logData, err := event.AWSLogs.Parse()
		if err != nil {
			log.Error("failed to parse event: %s", err.Error())
			return
		}

		for _, logEvent := range logData.LogEvents {
			_, err := fmt.Fprintln(w, logEvent.Message)
			if err != nil {
				log.Error("failed to push log event: %s", err)
			}
		}
	})
}
