package db

import (
	"context"
)

type Manager struct {
	client  Client
	channel chan Request
}

func NewManager(client Client) *Manager {
	return &Manager{client: client, channel: make(chan Request)}
}

func (m *Manager) Start(ctx context.Context) error {
	for req := range m.channel {
		result := m.client.ExecuteQuery(ctx, req.Query, req.Args...)
		if result.Error != nil {
			req.ResultChan <- Result{Error: result.Error}
			continue
		}

		// TODO: populate data from result

		req.ResultChan <- Result{}
	}

	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	close(m.channel)
	return nil
}
