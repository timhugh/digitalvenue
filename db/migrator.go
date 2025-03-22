package db

import "time"

type Status string

const (
	MigrationPending    Status = "pending"
	MigrationApplied    Status = "applied"
	MigrationFailed     Status = "failed"
	MigrationRolledBack Status = "rolledback"
)

type Migrator interface {
	MigrateAll() error
	MigrateOne(Migration) error
	Rollback(Migration) error
}

type Migration struct {
	ID     int
	Name   string
	Status Status
	Up     func() error
	Down   func() error
}

type MigrationStatus struct {
	Migration Migration
	AppliedAt time.Time
}
