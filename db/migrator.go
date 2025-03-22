package db

import (
	"context"
	"log"
	"time"
)

type Status string

const (
	MigrationPending    Status = "pending"
	MigrationApplied    Status = "applied"
	MigrationFailed     Status = "failed"
	MigrationRolledBack Status = "rolledback"
)

type Migration struct {
	Version int
	Name    string
	Status  Status
	Up      func(context.Context, Client) error
	Down    func(context.Context, Client) error
}

type MigrationRecord struct {
	Version   int       `json:"version" db:"version"`
	AppliedAt time.Time `json:"applied_at" db:"applied_at"`
	Status    Status    `json:"status" db:"status"`
}

type Migrator struct {
	Client     Client
	Migrations []Migration
}

func (m *Migrator) MigrateAll() error {
	for _, migration := range m.Migrations {
		log.Printf("Migration %d: %s", migration.Version, migration.Name)
	}
	return nil
}

func (m *Migrator) MigrateOne(migration Migration) error {
	return nil
}

func (m *Migrator) Rollback(migration Migration) error {
	return nil
}
