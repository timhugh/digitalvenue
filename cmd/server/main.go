package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/timhugh/digitalvenue/cmd/server/api"
)

type worker interface {
	Run(context.Context) error
	Stop(context.Context) error
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	workers := make([]worker, 0)
	workers = append(workers, api.NewWorker())

	var wg sync.WaitGroup
	for _, worker := range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := worker.Run(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "failed to start worker: %s\n", err)
				cancel()
			}
		}()
	}

	<-ctx.Done()
	fmt.Println("shutdown request received, stopping workers...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, worker := range workers {
		if err := worker.Stop(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "failed to stop worker: %s\n", err)
		}
	}

	wg.Wait()
	fmt.Println("Workers stopped. Goodbye!")

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
