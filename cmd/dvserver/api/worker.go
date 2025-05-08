package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/timhugh/digitalvenue/util/events"
	"net"
	"net/http"
	"time"
)

func NewWorker(bus *events.Bus) *Worker {
	return &Worker{
		bus: bus,
	}
}

type Worker struct {
	httpServer http.Server
	bus        *events.Bus
}

func (w *Worker) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	addRoutes(mux, w.bus)

	w.httpServer = http.Server{
		Addr:        net.JoinHostPort("0.0.0.0", "8080"),
		Handler:     mux,
		ReadTimeout: time.Duration(5 * time.Second),
	}

	if err := w.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server encountered error: %s\n", err)
	}

	return nil
}

func (w *Worker) Stop(ctx context.Context) error {
	if err := w.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server failed to shutdown: %s", err)
	}
	return nil
}
