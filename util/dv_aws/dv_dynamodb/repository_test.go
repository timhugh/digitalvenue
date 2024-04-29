package dv_dynamodb

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/timhugh/digitalvenue/util/core"
	"github.com/timhugh/digitalvenue/util/test"
	"os"
	"testing"
)

const tableName = "test-table"

type testItem struct {
	ID      string
	SomeKey string
}

func testItemAttributes() map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"ID":      &types.AttributeValueMemberS{Value: "ID"},
		"SomeKey": &types.AttributeValueMemberS{Value: "SomeValue"},
		"Type":    &types.AttributeValueMemberS{Value: "Object"},
	}
}

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

func TestRepository_get_success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	getItemInputCaptor := mock.Captor[*dynamodb.GetItemInput]()
	mock.When(client.GetItem(mock.AnyContext(), getItemInputCaptor.Capture())).
		ThenReturn(&dynamodb.GetItemOutput{Item: testItemAttributes()}, nil)

	var item testItem
	err := repo.get("Object", map[string]string{"ID": "ID"}, &item)
	if err != nil {
		t.Error(err)
	}

	expectedItem := testItem{
		ID:      "ID",
		SomeKey: "SomeValue",
	}
	if err := test.Diff(expectedItem, item); err != nil {
		t.Error(err)
	}

	expectedInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: "ID"},
		},
	}
	if err := test.Diff(expectedInput, getItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_get_ClientError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")
	mock.When(client.GetItem(mock.AnyContext(), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(nil, inducedError)

	err := repo.get("Object", map[string]string{}, nil)
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}

func TestRepository_get_IncorrectObjectType(t *testing.T) {
	repo, client := initRepositoryTest(t)

	mock.When(client.GetItem(mock.AnyContext(), mock.Any[*dynamodb.GetItemInput]())).
		ThenReturn(&dynamodb.GetItemOutput{Item: testItemAttributes()}, nil)

	var item testItem
	err := repo.get("SomeOtherObject", map[string]string{"ID": "ID"}, &item)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRepository_put_success(t *testing.T) {
	repo, client := initRepositoryTest(t)

	putItemInputCaptor := mock.Captor[*dynamodb.PutItemInput]()
	mock.When(client.PutItem(mock.Any[context.Context](), putItemInputCaptor.Capture())).
		ThenReturn(nil, nil)

	object := testItem{
		ID:      "ID",
		SomeKey: "SomeValue",
	}

	err := repo.put("Object", object)
	if err != nil {
		t.Error(err)
	}

	expectedInput := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      testItemAttributes(),
	}
	if err := test.Diff(expectedInput, putItemInputCaptor.Last()); err != nil {
		t.Error(err)
	}
}

func TestRepository_put_ClientError(t *testing.T) {
	repo, client := initRepositoryTest(t)

	inducedError := errors.New("induced error")
	mock.When(client.PutItem(mock.AnyContext(), mock.Any[*dynamodb.PutItemInput]())).
		ThenReturn(nil, inducedError)

	err := repo.put("Object", nil)
	if !errors.Is(err, inducedError) {
		t.Errorf("expected error %v, got %v", inducedError, err)
	}
}
