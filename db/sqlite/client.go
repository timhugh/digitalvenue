package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/timhugh/dv-go/db"
)

type Client struct {
	db *sql.DB
}

func New(databaseFile string) (*Client, error) {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		return nil, err
	}

	// TODO: enable WAL

	return &Client{db: db}, nil
}

func (c *Client) ExecuteQuery(ctx context.Context, query string, args ...any) db.Result {
	_, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return db.Result{Error: err}
	}

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return db.Result{Error: err}
	}
	defer rows.Close()

	return db.Result{}
}
