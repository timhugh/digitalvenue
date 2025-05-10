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

func main() {
	db_uri := os.Getenv("DV_DATABASE_URL")
	
	db_name, err := extractDatabaseName(db_uri)
	panicErr(err)
	pg_db_uri := strings.Replace(db_uri, db_name, "postgres", 1)
	pgdb, err := db.New(pg_db_uri)
	panicErr(err)
	defer pgdb.Close()

	res, err := pgdb.Exec("SELECT 1 FROM pg_database WHERE datname=$1", db_name)
	panicErr(err)

	rows, err := res.RowsAffected()
	if rows == 0 {
		fmt.Printf("database %s does not exist yet, creating", db_name)
		_, err := pgdb.Exec("CREATE DATABASE $1", db_name)
		panicErr(err)
	}

	db, err := db.New(db_uri)
	panicErr(err)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)
	err = goose.SetDialect("postgres")
	panicErr(err)
	err = goose.Up(db, "migrations")
	panicErr(err)
}
