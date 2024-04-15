package dv_dynamodb

import (
	"errors"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/core"
	"os"
	"testing"
)

const tableName = "test-table"

func initRepositoryTest(t *testing.T) (*Repository, Client) {
	mock.SetUp(t)

	err := os.Setenv("CORE_DATA_TABLE_NAME", tableName)
	if err != nil {
		t.Fatal(err)
	}

	client := mock.Mock[Client]()

	repo, err := NewRepository(client)
	if err != nil {
		t.Fatal(err)
	}

	return repo, client
}

func TestNewRepository_RequiresTableName(t *testing.T) {
	err := os.Unsetenv("CORE_DATA_TABLE_NAME")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRepository(nil)
	if err == nil {
		t.Error("expected error, got nil")
	}

	var keyError core.MissingEnvError
	ok := errors.As(err, &keyError)
	if !ok {
		t.Errorf("expected error to be of type MissingEnvError, got %T", err)
	}

	err = os.Setenv("CORE_DATA_TABLE_NAME", "test-table")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewRepository(nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
