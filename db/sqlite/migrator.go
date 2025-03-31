package sqlite

import (
	"context"
	"embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/timhugh/digitalvenue/db"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func NewMigrator(client db.Client) (*db.Migrator, error) {
	migrations, err := loadMigrations()
	if err != nil {
		return nil, err
	}

	result := client.ExecuteQuery(context.Background(), "PRAGMA foreign_keys = ON")
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to enable foreign keys: %w", result.Error)
	}

	result = client.ExecuteQuery(context.Background(), "PRAGMA journal_mode = WAL")
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to set journal mode to WAL: %w", result.Error)
	}

	versionsMigration := migrations[0]
	err = versionsMigration.Up(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("Failed to create versions table: %w", err)
	}

	return &db.Migrator{
		Client:     client,
		Migrations: migrations[1:],
	}, nil
}

func loadMigrations() ([]db.Migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, err
	}
	var migrations []db.Migration
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		migration, err := loadMigration(entry.Name())
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}
	return migrations, nil
}

func loadMigration(filename string) (db.Migration, error) {
	name, version, err := getNameAndVersion(filename)
	if err != nil {
		return db.Migration{}, err
	}

	content, err := migrationsFS.ReadFile("migrations/" + filename)
	if err != nil {
		return db.Migration{}, fmt.Errorf("failed to read migration file %s: %w", filename, err)
	}

	sections := strings.Split(string(content), "-- Down")
	upSQL := strings.TrimSpace(sections[0])
	upSQL = strings.TrimPrefix(upSQL, "-- Up")
	upSQL = strings.TrimSpace(upSQL)
	downSQL := strings.TrimSpace(sections[1])

	up := func(ctx context.Context, client db.Client) error {
		result := client.ExecuteQuery(ctx, upSQL)
		return result.Error
	}
	down := func(ctx context.Context, client db.Client) error {
		result := client.ExecuteQuery(ctx, downSQL)
		return result.Error
	}

	return db.Migration{
		Name:    name,
		Version: version,
		Up:      up,
		Down:    down,
	}, nil
}

func getNameAndVersion(filename string) (name string, version int, err error) {
	parts := strings.Split(filename, "-")
	if len(parts) < 2 {
		return "", 0, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	name = strings.Join(parts[1:], " ")
	name = strings.TrimSuffix(name, ".sql")

	version, err = strconv.Atoi(parts[0])
	if err != nil {
		return "", 0, fmt.Errorf("invalid migration version: %s", filename)
	}

	return name, version, nil
}
