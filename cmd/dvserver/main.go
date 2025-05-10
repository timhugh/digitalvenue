package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/timhugh/digitalvenue/cmd/dvserver/api"
	"github.com/timhugh/digitalvenue/cmd/dvserver/test"
	"github.com/timhugh/digitalvenue/util/events"
)

type worker interface {
	Run(context.Context) error
	Stop(context.Context) error
}

const workerShutdownTimeoutSeconds = 10

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	eventsBus := events.NewBus()

	workers := make([]worker, 0)
	workers = append(workers, api.NewWorker(eventsBus))
	workers = append(workers, test.NewWorker(eventsBus))

	fmt.Printf("booting %d workers...\n", len(workers))
	var wg sync.WaitGroup
	for i, worker := range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			fmt.Printf("starting worker #%d...\n", i)
			if err := worker.Run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "failed to start worker #%d: %s\n", i, err)
				cancel()
			}
		}()
	}

	fmt.Println("start request sent to all workers")
	fmt.Println("main loop waiting for stop")
	<-ctx.Done()
	fmt.Printf("shutdown request received, stopping workers within %d seconds\n", workerShutdownTimeoutSeconds)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), workerShutdownTimeoutSeconds*time.Second)
	defer cancel()
	for _, worker := range workers {
		if err := worker.Stop(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to stop worker: %s\n", err)
		}
	}

	wg.Wait()
	fmt.Println("workers stopped. Goodbye!")

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
