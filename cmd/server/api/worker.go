package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func NewWorker() *Worker {
	return &Worker{}
}

type Worker struct {
	httpServer http.Server
}

func (w *Worker) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	addRoutes(mux)

	w.httpServer = http.Server{
		Addr:        net.JoinHostPort("0.0.0.0", "8080"),
		Handler:     mux,
		ReadTimeout: time.Duration(5 * time.Second),
	}

	if err := w.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
