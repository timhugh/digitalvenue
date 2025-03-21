package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func createTestDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE users (uuid TEXT PRIMARY KEY, name TEXT NOT NULL)")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestExecuteQuery(t *testing.T) {
	type User struct {
		UUID string `json:"uuid" db:"uuid"`
		Name string `json:"name" db:"name"`
	}

	db, err := createTestDB()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	johnDoe := User{UUID: uuid.NewString(), Name: "John Doe"}
	janeDoe := User{UUID: uuid.NewString(), Name: "Jane Doe"}

	_, err = db.Exec("INSERT INTO users (uuid, name) VALUES (?, ?)", johnDoe.UUID, johnDoe.Name)
	if err != nil {
		t.Fatalf("Failed to insert user %s: %v", johnDoe.Name, err)
	}
	_, err = db.Exec("INSERT INTO users (uuid, name) VALUES (?, ?)", janeDoe.UUID, janeDoe.Name)
	if err != nil {
		t.Fatalf("Failed to insert user %s: %v", janeDoe.Name, err)
	}

	client := Client{db: db}

	t.Run("Simple retrieval", func(t *testing.T) {
		result := client.ExecuteQuery(context.Background(), "SELECT * FROM users")
		if result.Error != nil {
			t.Fatalf("Expected no error, got %v", result.Error)
		}

		var data []User
		err = result.Unwrap(&data)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(data) != 2 {
			t.Fatalf("Expected 2 users, got %d", len(data))
		}
		if diff := deep.Equal(data[0], johnDoe); diff != nil {
			t.Errorf("Expected user %#v, got diff %#v", johnDoe, diff)
		}
		if diff := deep.Equal(data[1], janeDoe); diff != nil {
			t.Errorf("Expected user %#v, got diff %#v", janeDoe, diff)
		}
	})

	t.Run("Retrieve with args", func(t *testing.T) {
		result := client.ExecuteQuery(context.Background(), "SELECT * FROM users WHERE uuid = ?", johnDoe.UUID)
		if result.Error != nil {
			t.Fatalf("Expected no error, got %v", result.Error)
		}

		var data User
		err = result.Unwrap(&data)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if diff := deep.Equal(data, johnDoe); diff != nil {
			t.Errorf("Expected user %#v, got diff %#v", johnDoe, diff)
		}
	})

	t.Run("Simple insert", func(t *testing.T) {
		user := User{UUID: uuid.NewString(), Name: "Alice"}
		result := client.ExecuteQuery(context.Background(), "INSERT INTO users (uuid, name) VALUES (?, ?)", user.UUID, user.Name)
		if result.Error != nil {
			t.Fatalf("Expected no error, got %v", result.Error)
		}
		row := db.QueryRow("SELECT * FROM users WHERE uuid = ?", user.UUID)
		if row.Err() != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		var userResult User
		if err := row.Scan(&userResult.UUID, &userResult.Name); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if diff := deep.Equal(userResult, user); diff != nil {
			t.Errorf("Expected user %#v, got diff %#v", user, diff)
		}
	})

	t.Run("Simple update", func(t *testing.T) {
		result := client.ExecuteQuery(context.Background(), "UPDATE users SET name = ? WHERE uuid = ?", "Jonathan", johnDoe.UUID)
		if result.Error != nil {
			t.Fatalf("Expected no error, got %v", result.Error)
		}
		row := db.QueryRow("SELECT * FROM users WHERE uuid = ?", johnDoe.UUID)
		if row.Err() != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		var userResult User
		if err := row.Scan(&userResult.UUID, &userResult.Name); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expectedUser := User{UUID: johnDoe.UUID, Name: "Jonathan"}
		if diff := deep.Equal(userResult, expectedUser); diff != nil {
			t.Errorf("Expected user %#v, got diff %#v", expectedUser, diff)
		}
	})

	t.Run("Simple delete", func(t *testing.T) {
		result := client.ExecuteQuery(context.Background(), "DELETE FROM users WHERE uuid = ?", johnDoe.UUID)
		if result.Error != nil {
			t.Fatalf("Expected no error, got %v", result.Error)
		}
		row := db.QueryRow("SELECT * FROM users WHERE uuid = ?", johnDoe.UUID)
		if row.Err() != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		var userResult User
		if err := row.Scan(&userResult.UUID, &userResult.Name); err != sql.ErrNoRows {
			t.Errorf("expected no rows to exist, got user %v", userResult)
		}
	})
}
