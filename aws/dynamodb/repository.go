package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
	"strings"
)

const coreDataTableName = "CORE_DATA_TABLE_NAME"
const itemTypeKey = "Type"

type Repository struct {
	client    Client
	tableName string
}

func NewRepository(client Client) (*Repository, error) {
	tableName, err := core.RequireEnv(coreDataTableName)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (repo *Repository) getItem(itemType string, getItemInput *dynamodb.GetItemInput, out interface{}) error {
	getItemOutput, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return errors.Wrap(err, "error retrieving item from dynamodb")
	}

	if getItemOutput.Item == nil {
		return errors.New("item not found")
	}

	retrievedItemType := getItemOutput.Item[itemTypeKey].(*types.AttributeValueMemberS).Value
	if retrievedItemType != itemType {
		return errors.New("retrieved item is not a " + itemType)
	}

	err = attributevalue.UnmarshalMap(getItemOutput.Item, &out)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling dynamodb get item output")
	}

	return nil
}

func removeIDPrefix(prefixedID string) (string, error) {
	parts := strings.Split(prefixedID, "#")
	if len(parts) <= 1 {
		return "", errors.New("invalid prefixed ID")
	}

	return parts[1], nil
}
