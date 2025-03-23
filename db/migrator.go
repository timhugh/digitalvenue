package db

import (
	"context"
	"fmt"
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
	Name      string    `json:"name" db:"name"`
	AppliedAt time.Time `json:"applied_at" db:"applied_at"`
	Status    Status    `json:"status" db:"status"`
}

type VersionRepository interface {
	All(ctx context.Context) ([]MigrationRecord, error)
}

type Migrator struct {
	Client     Client
	Versions   VersionRepository
	Migrations []Migration
}

func isMigrationApplied(migration Migration, migrationRecords []MigrationRecord) bool {
	for _, migrationRecord := range migrationRecords {
		if migration.Version == migrationRecord.Version && migrationRecord.Status == MigrationApplied {
			return true
		}
	}
	return false
}

func (m *Migrator) MigrateAll(ctx context.Context) error {
	migrationRecords, err := m.Versions.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch migration versions: %w", err)
	}
	for _, migration := range m.Migrations {
		if isMigrationApplied(migration, migrationRecords) {
			continue
		}

		if err := m.MigrateOne(ctx, migration); err != nil {
			return fmt.Errorf("migration failed %s: %w", migration.Name, err)
		}
	}
	return nil
}

func (m *Migrator) MigrateOne(ctx context.Context, migration Migration) error {
	err := migration.Up(ctx, m.Client)
	if err != nil {
		return fmt.Errorf("migration failed %s: %w", migration.Name, err)
	}
}

func (m *Migrator) Rollback(ctx context.Context, migration Migration) error {
	return migration.Down(ctx, m.Client)
}
