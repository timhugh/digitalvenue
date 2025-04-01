package test

import "context"

func NewWorker() *Worker {
	return &Worker{}
}

type Worker struct {
}

func (w *Worker) Run(ctx context.Context) {
}

func (w *Worker) Stop() {
}
