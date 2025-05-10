package main

import (
	"embed"
	"fmt"
	"net/url"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/timhugh/digitalvenue/db"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func extractDatabaseName(dbURI *url.URL) (string, error) {
	path := dbURI.Path
	if path == "" || path == "/" {
		return "", fmt.Errorf("database URI is missing database name")
	}

	return strings.TrimPrefix(path, "/"), nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createDbIfNotExists(dbURI string) error {
	URI, err := url.Parse(dbURI)
	if err != nil {
		return fmt.Errorf("failed to parse database URI: %w", err)
	}

	dbName, err := extractDatabaseName(URI)
	if err != nil {
		return fmt.Errorf("failed to extract db name from connection string: %w", err)
	}

	URI.Path = "/postgres"

	pgdb, err := db.New(URI.String())
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer pgdb.Close()

	res, err := pgdb.Exec("SELECT 1 FROM pg_database WHERE datname=$1", dbName)
	if err != nil {
		return fmt.Errorf("failed to query for existing database %s: %w", dbName, err)
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		fmt.Printf("database %s does not exist yet, creating\n", dbName)
		if _, err := pgdb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
			return fmt.Errorf("failed to create database %s: %w", dbName, err)
		}
	}

	return nil
}

func main() {
	dbURI := os.Getenv("DV_DATABASE_URL")
	panicIfErr(createDbIfNotExists(dbURI))

	db, err := db.New(dbURI)
	panicIfErr(err)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)
	panicIfErr(goose.SetDialect("postgres"))
	fmt.Println("executing migrations")
	panicIfErr(goose.Up(db, "migrations"))
}
