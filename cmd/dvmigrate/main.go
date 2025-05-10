package main

import (
	"embed"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/timhugh/digitalvenue/db"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	db_uri := os.Getenv("DV_DATABASE_URL")
	db, err := db.New(db_uri)
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
