package sqlite

import "github.com/timhugh/dv-go/db"

type Migrator struct {
	client *Client
}

func (m *Migrator) MigrateAll() error {
	return nil
}

func (m *Migrator) MigrateOne(migration db.Migration) error {
	return nil
}

func (m *Migrator) Rollback(migration db.Migration) error {
	return nil
}

func (m *Migrator) LoadMigrations() error {
	return nil
}
