package db

import (
	"context"
	"fmt"
	"time"

	"github.com/timhugh/digitalvenue/util/optional"
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
	Up      func(context.Context, SqlClient) error
	Down    func(context.Context, SqlClient) error
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
	Client     SqlClient
	Versions   VersionRepository
	Migrations []Migration
}

func (m *Migrator) MigrateAll(ctx context.Context) error {
	migrationRecords, err := m.Versions.All(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch migration versions: %w", err)
	}

	for _, migration := range m.Migrations {
		var migrationRecord = optional.Empty[MigrationRecord]()
		for _, record := range migrationRecords {
			if migration.Version == record.Version {
				migrationRecord.Set(record)
			}
		}

		if err := m.MigrateOne(ctx, migration, migrationRecord); err != nil {
			return fmt.Errorf("migration failed %s: %w", migration.Name, err)
		}
	}
	return nil
}

func (m *Migrator) MigrateOne(ctx context.Context, migration Migration, record optional.Optional[MigrationRecord]) error {
	err := migration.Up(ctx, m.Client)
	if err != nil {
		return fmt.Errorf("migration failed %s: %w", migration.Name, err)
	}

	return nil
}

func (m *Migrator) Rollback(ctx context.Context, migration Migration) error {
	return migration.Down(ctx, m.Client)
}
