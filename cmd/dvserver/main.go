package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/timhugh/digitalvenue/cmd/dvserver/api"
	"github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/db/queries"
	"github.com/timhugh/digitalvenue/util/events"
	"github.com/timhugh/digitalvenue/util/logger"
)

type worker interface {
	Run(context.Context) error
	Stop(context.Context) error
}

const workerShutdownTimeoutSeconds = 10

func run(ctx context.Context, out io.Writer, errOut io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	logLevel := os.Getenv("DV_LOG_LEVEL")
	log := logger.New(out, errOut, logLevel)

	eventsBus := events.NewBus()
	dbURI := os.Getenv("DV_DATABASE_URL")
	db, err := db.New(dbURI)
	if err != nil {
		return fmt.Errorf("failed to connect to database %s: %w", dbURI, err)
	}
	queries := queries.New(db)

	workers := make([]worker, 0)
	workers = append(workers, api.NewWorker(eventsBus, queries, log))

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
	if err := run(ctx, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
