package main

import (
	"context"

	"github.com/timhugh/digitalvenue/db"
	"github.com/timhugh/digitalvenue/db/sqlite"
)

func main() {
	var db db.Client
	db, err := sqlite.NewClient("test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrator, err := sqlite.NewMigrator(db)
	if err != nil {
		panic(err)
	}

	err = migrator.MigrateAll(context.Background())

	if err != nil {
		panic(err)
	}
}
