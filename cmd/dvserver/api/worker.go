package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/timhugh/digitalvenue/util/events"
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

const serverAddress = "0.0.0.0"
const serverPort = "8080"

func (w *Worker) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	addRoutes(mux, w.bus)

	w.httpServer = http.Server{
		Addr:        net.JoinHostPort(serverAddress, serverPort),
		Handler:     mux,
		ReadTimeout: time.Duration(5 * time.Second),
	}

	fmt.Printf("starting web server at %s:%s\n", serverAddress, serverPort)
	if err := w.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server encountered error: %s", err)
	}

	return nil
}

func (w *Worker) Stop(ctx context.Context) error {
	if err := w.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server failed to shutdown: %s", err)
	}
	return nil
}
