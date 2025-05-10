package db

import "database/sql"

func New(connection string) (*sql.DB, error) {
	return sql.Open("postgres", connection)
}
