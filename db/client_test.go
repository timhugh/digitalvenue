package db_test

import (
	"errors"
	"testing"

	"github.com/timhugh/digitalvenue/db"
)

func TestResultUnwrap(t *testing.T) {
	type User struct {
		ID   int    `json:"id" db:"id"`
		Name string `json:"name" db:"name"`
	}

	t.Run("Unwraps struct data", func(t *testing.T) {
		result := db.Result{
			Data: []map[string]any{
				{"id": 1, "name": "John"},
			},
			Error: nil,
		}
		var user User
		err := result.Unwrap(&user)
		if err != nil {
			t.Fatalf("failed to unwrap user: %v", err)
		}
		if user.ID != 1 {
			t.Fatalf("expected user ID 1, got %d", user.ID)
		}
		if user.Name != "John" {
			t.Fatalf("expected user name John, got %s", user.Name)
		}
	})

	t.Run("Unwraps slice data", func(t *testing.T) {
		result := db.Result{
			Data: []map[string]any{
				{"id": 1, "name": "John"},
				{"id": 2, "name": "Jane"},
			},
			Error: nil,
		}
		var users []User
		err := result.Unwrap(&users)
		if err != nil {
			t.Fatalf("failed to unwrap users: %v", err)
		}
		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}
		if users[0].ID != 1 {
			t.Fatalf("expected user ID 1, got %d", users[0].ID)
		}
		if users[0].Name != "John" {
			t.Fatalf("expected user name John, got %s", users[0].Name)
		}
		if users[1].ID != 2 {
			t.Fatalf("expected user ID 2, got %d", users[1].ID)
		}
		if users[1].Name != "Jane" {
			t.Fatalf("expected user name Jane, got %s", users[1].Name)
		}
	})

	t.Run("Errors on empty data", func(t *testing.T) {
		result := db.Result{
			Data:  []map[string]any{},
			Error: nil,
		}
		var output string
		err := result.Unwrap(&output)
		if !errors.Is(err, db.ErrNoResults) {
			t.Fatalf("expected ErrNoResults, got %v", err)
		}
	})

	t.Run("Errors on invalid output type", func(t *testing.T) {
		result := db.Result{
			Data:  []map[string]any{{"id": 1, "name": "John"}},
			Error: nil,
		}
		var output string
		err := result.Unwrap(&output)
		if !errors.Is(err, db.ErrMappingFailed) {
			t.Fatalf("expected ErrMappingFailed, got %v", err)
		}
	})
}
