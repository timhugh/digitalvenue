package main

import (
	"github.com/timhugh/dv-go/db"
	"github.com/timhugh/dv-go/db/sqlite"
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

	err = migrator.MigrateAll()

	if err != nil {
		panic(err)
	}
}
