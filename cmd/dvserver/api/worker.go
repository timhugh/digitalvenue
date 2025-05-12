package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/timhugh/digitalvenue/common"
	"github.com/timhugh/digitalvenue/db/queries"
	"github.com/timhugh/digitalvenue/util/events"
)

func NewWorker(bus *events.Bus, queries *queries.Queries, log common.Logger) *Worker {
	return &Worker{
		log:     log,
		bus:     bus,
		queries: queries,
	}
}

type Worker struct {
	httpServer http.Server
	log        common.Logger
	bus        *events.Bus
	queries    *queries.Queries
}

const serverAddress = "0.0.0.0"
const serverPort = "8080"

func (w *Worker) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	addRoutes(mux, w.log, w.bus, w.queries)

	stack := LoggingMiddleware(w.log)(mux)

	w.httpServer = http.Server{
		Addr:        net.JoinHostPort(serverAddress, serverPort),
		Handler:     stack,
		ReadTimeout: time.Duration(5 * time.Second),
	}

	w.log.Info("starting web server at %s:%s", serverAddress, serverPort)
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
