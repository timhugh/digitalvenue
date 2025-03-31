package sqlite

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/timhugh/digitalvenue/db"
)

const versionsTable = "versions"

var versionsColumns = []string{"version", "name", "applied_at", "status"}

type VersionRepository struct {
	db      db.Client
	builder squirrel.StatementBuilderType
}

func NewVersionRepository(db db.Client) *VersionRepository {
	return &VersionRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *VersionRepository) All(ctx context.Context) ([]db.MigrationRecord, error) {
	query, args, err := r.builder.Select(versionsColumns...).From(versionsTable).ToSql()
	if err != nil {
		return nil, fmt.Errorf("Failed to build SQL query: %w", err)
	}

	result := r.db.ExecuteQuery(ctx, query, args...)
	if result.Error != nil {
		return nil, fmt.Errorf("Failed to execute query: %w", result.Error)
	}

	var records []db.MigrationRecord
	err = result.Unwrap(&records)
	if err != nil {
		return nil, fmt.Errorf("Failed to unwrap results: %w", err)
	}

	return records, nil
}

func (r *VersionRepository) Upsert(ctx context.Context, record db.MigrationRecord) error {
	existingRecord, err := r.GetByVersion(ctx, record.Version)

	if err != nil {
		return fmt.Errorf("Failed to build SQL query: %w", err)
	}

	result := r.db.ExecuteQuery(ctx, query, args...)
	if result.Error != nil {
		return fmt.Errorf("Failed to execute query: %w", result.Error)
	}

	return nil
}

func (r *VersionRepository) GetByVersion(ctx context.Context, version int) (db.MigrationRecord, error) {
	query, args, err := r.builder.Select(versionsColumns...).From(versionsTable).Where(squirrel.Eq{"version": version}).ToSql()
	if err != nil {
		return db.MigrationRecord{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}

	result := r.db.ExecuteQuery(ctx, query, args...)
	if result.Error != nil {
		return db.MigrationRecord{}, fmt.Errorf("Failed to execute query: %w", result.Error)
	}

	var record db.MigrationRecord
	err = result.Unwrap(&record)
	if err != nil {
		return db.MigrationRecord{}, fmt.Errorf("Failed to unwrap results: %w", err)
	}

	return record, nil
}
