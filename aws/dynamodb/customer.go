package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/core"
)

type CustomerRepository struct {
	client      Client
	idGenerator core.IDGenerator
	tableName   string
}

func NewCustomerRepository(client Client) (*CustomerRepository, error) {
	tableName, err := core.RequireEnv(CustomersTableNameKey)
	if err != nil {
		return nil, err
	}

	return &CustomerRepository{
		client:      client,
		idGenerator: core.NewIDGenerator(),
		tableName:   tableName,
	}, nil
}

func (repo *CustomerRepository) PutCustomer(customer core.Customer) (string, error) {
	var customerID string
	if customer.CustomerID == "" {
		customerID = repo.idGenerator.GenerateID()
	} else {
		customerID = customer.CustomerID
	}

	var meta = make(map[string]types.AttributeValue)
	for key, value := range customer.Meta {
		meta[key] = &types.AttributeValueMemberS{Value: value}
	}

	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			CustomerID: &types.AttributeValueMemberS{Value: customerID},
			FirstName:  &types.AttributeValueMemberS{Value: customer.FirstName},
			LastName:   &types.AttributeValueMemberS{Value: customer.LastName},
			Email:      &types.AttributeValueMemberS{Value: customer.Email},
			Phone:      &types.AttributeValueMemberS{Value: customer.Phone},
			Meta:       &types.AttributeValueMemberM{Value: meta},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return "", err
	}

	return customerID, nil
}
