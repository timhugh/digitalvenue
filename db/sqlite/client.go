package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/dv-go/db"
)

type Client struct {
	db *sqlx.DB
}

func New(databaseFile string) (*Client, error) {
	dbHandle, err := sqlx.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w %w", db.ErrInitFailed, err)
	}

	_, err = dbHandle.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to set WAL mode: %w %w", db.ErrInitFailed, err)
	}

	return &Client{db: dbHandle}, nil
}

func (c *Client) ExecuteQuery(ctx context.Context, query string, args ...any) db.Result {
	rows, err := c.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return db.Result{Error: fmt.Errorf("failed to execute query: %w %w", db.ErrQueryFailed, err)}
	}
	defer rows.Close()

	var data []map[string]any
	for rows.Next() {
		row := make(map[string]any)
		if err := rows.MapScan(row); err != nil {
			return db.Result{Error: fmt.Errorf("failed to scan row: %w %w", db.ErrQueryFailed, err)}
		}
		data = append(data, row)
	}

	return db.Result{Data: data}
}
