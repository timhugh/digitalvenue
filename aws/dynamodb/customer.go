package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/aws/dynamodb/square"
	"github.com/timhugh/digitalvenue/core"
	"os"
)

type CustomerRepositoryConfig struct {
	TableName string
}

func NewCustomerRepositoryConfig() CustomerRepositoryConfig {
	return CustomerRepositoryConfig{
		TableName: os.Getenv(CustomersTableName),
	}
}

type customerRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewCustomerRepository(config CustomerRepositoryConfig, client *dynamodb.Client) core.CustomerRepository {
	return customerRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo customerRepository) Put(customer core.Customer) (string, error) {
	var customerID string
	if customer.CustomerID == "" {
		customerID = core.GenerateID()
	} else {
		customerID = customer.CustomerID
	}

	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			CustomerID: &types.AttributeValueMemberS{Value: customerID},
			FirstName:  &types.AttributeValueMemberS{Value: customer.FirstName},
			LastName:   &types.AttributeValueMemberS{Value: customer.LastName},
			Email:      &types.AttributeValueMemberS{Value: customer.Email},
			Phone:      &types.AttributeValueMemberS{Value: customer.Phone},
			Meta: &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
				square.SquareCustomerID: &types.AttributeValueMemberS{Value: customer.Meta.SquareCustomerID},
			}},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return "", err
	}

	return customerID, nil
}
