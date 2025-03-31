package sqlite

import (
	"context"
	"regexp"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/timhugh/digitalvenue/db"
)

type mockClient struct {
	ReceivedContext context.Context
	ReceivedQuery   string
	ReceivedArgs    []any

	ReturnedResult db.Result
}

func (m *mockClient) ExecuteQuery(ctx context.Context, query string, args ...any) db.Result {
	m.ReceivedContext = ctx
	m.ReceivedQuery = query
	m.ReceivedArgs = args
	return m.ReturnedResult
}

func (m *mockClient) Close() error {
	return nil
}

var multiSpaceRegex = regexp.MustCompile(`\s+`)

func trimQuery(query string) string {
	return multiSpaceRegex.ReplaceAllString(query, " ")
}

func TestGetNameAndVersion(t *testing.T) {
	filename := "0000-create-versions.sql"

	name, version, err := getNameAndVersion(filename)
	if err != nil {
		t.Fatalf("Error getting name and version: %v", err)
	}

	if name != "create versions" {
		t.Errorf("Unexpected migration name: %s", name)
	}
	if version != 0 {
		t.Errorf("Unexpected migration version: %d", version)
	}
}

func TestLoadMigrations(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("Error loading migrations: %v", err)
	}

	versionMigration := migrations[0]

	t.Run("Up migration query", func(t *testing.T) {
		mockClient := &mockClient{}
		ctx := context.Background()
		if err := versionMigration.Up(ctx, mockClient); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedQuery := `CREATE TABLE IF NOT EXISTS versions (
      version INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
      status TEXT NOT NULL DEFAULT 'pending'
    );`

		if expected, actual := trimQuery(expectedQuery), trimQuery(mockClient.ReceivedQuery); expected != actual {
			t.Errorf("Query does not match expected: \n%v", diff.LineDiff(expected, actual))
		}
		if mockClient.ReceivedArgs != nil {
			t.Errorf("Unexpected args: %v", mockClient.ReceivedArgs)
		}
	})

	t.Run("Down migration query", func(t *testing.T) {
		mockClient := &mockClient{}
		ctx := context.Background()
		if err := versionMigration.Down(ctx, mockClient); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedQuery := "DROP TABLE versions;"
		if mockClient.ReceivedQuery != expectedQuery {
			t.Errorf("Expected query:\n %s\ngot:\n%s", expectedQuery, mockClient.ReceivedQuery)
		}
		if mockClient.ReceivedArgs != nil {
			t.Errorf("Unexpected args: %v", mockClient.ReceivedArgs)
		}
	})
}
