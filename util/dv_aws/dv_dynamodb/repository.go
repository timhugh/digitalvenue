package dv_dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"strings"
)

const coreDataTableNameKey = "CORE_DATA_TABLE_NAME"
const itemTypeKey = "Type"

type ItemNotFoundException struct {
	error
}

type Repository struct {
	client    Client
	tableName string
}

func NewRepository(client Client) (*Repository, error) {
	tableName, err := core.RequireEnv(coreDataTableNameKey)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (repo *Repository) get(itemType string, key map[string]string, out interface{}) error {
	keyAttrs := make(map[string]types.AttributeValue)
	for k, v := range key {
		keyAttrs[k] = &types.AttributeValueMemberS{Value: v}
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key:       keyAttrs,
	}

	getItemOutput, err := repo.client.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return errors.Wrap(err, "error retrieving item from dynamodb")
	}

	if getItemOutput.Item == nil {
		return ItemNotFoundException{}
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

func (repo *Repository) put(itemType string, itemDTO interface{}) error {
	attrs, err := attributevalue.MarshalMap(itemDTO)
	if err != nil {
		return errors.Wrap(err, "error marshalling dynamodb put item input")
	}
	attrs["Type"] = &types.AttributeValueMemberS{Value: itemType}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      attrs,
	}

	_, err = repo.client.PutItem(context.TODO(), input)
	if err != nil {
		return errors.Wrap(err, "error putting dynamodb item")
	}

	return nil
}

func UnprefixID(prefixedID string) (string, error) {
	parts := strings.Split(prefixedID, "#")
	if len(parts) <= 1 {
		return "", errors.New("invalid prefixed ID")
	}

	return parts[1], nil
}

func PrefixID(prefix string, id string) string {
	return prefix + "#" + id
}
