package test

import (
	"context"
	"github.com/timhugh/digitalvenue/util/events"
	"log"
)

func NewWorker(bus *events.Bus) *Worker {
	return &Worker{
		bus: bus,
	}
}

type Worker struct {
	bus *events.Bus
}

func (w *Worker) Run(ctx context.Context) error {
	sub := events.Subscribe[string](w.bus, 10)
	for {
		select {
		case event := <-sub.Events:
			log.Println(event)
		case <-ctx.Done():
			return nil
		}
	}
}

func (w *Worker) Stop(ctx context.Context) error {
	return nil
}
