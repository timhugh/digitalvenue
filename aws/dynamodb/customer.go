package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/core"
)

type customer struct {
	PK         string
	SK         string
	Type       string
	CustomerID string
	Name       string
	Email      string
	Phone      string
	Meta       map[string]string
}

func (repo *Repository) GetCustomer(tenantID string, customerID string) (*core.Customer, error) {
	tenantKey := "Tenant#" + tenantID
	customerKey := "Customer#" + customerID
	input := &dynamodb.GetItemInput{
		TableName: aws.String(repo.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: tenantKey},
			"SK": &types.AttributeValueMemberS{Value: customerKey},
		},
	}

	item := customer{}
	err := repo.getItem("Customer", input, &item)
	if err != nil {
		return nil, err
	}

	var meta map[string]string
	if item.Meta != nil {
		meta = make(map[string]string)
		for k, v := range item.Meta {
			meta[k] = v
		}
	}

	return &core.Customer{
		TenantID: tenantID,
		ID:       customerID,
		Name:     item.Name,
		Email:    item.Email,
		Phone:    item.Phone,
		Meta:     meta,
	}, nil
}

func (repo *Repository) PutCustomer(customer core.Customer) error {
	tenantKey := "Tenant#" + customer.TenantID
	customerKey := "Customer#" + customer.ID

	var meta map[string]types.AttributeValue
	if customer.Meta != nil {
		meta = make(map[string]types.AttributeValue)
		for k, v := range customer.Meta {
			meta[k] = &types.AttributeValueMemberS{Value: v}
		}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item: map[string]types.AttributeValue{
			"PK":         &types.AttributeValueMemberS{Value: tenantKey},
			"SK":         &types.AttributeValueMemberS{Value: customerKey},
			"Type":       &types.AttributeValueMemberS{Value: "Customer"},
			"CustomerID": &types.AttributeValueMemberS{Value: customer.ID},
			"Name":       &types.AttributeValueMemberS{Value: customer.Name},
			"Email":      &types.AttributeValueMemberS{Value: customer.Email},
			"Phone":      &types.AttributeValueMemberS{Value: customer.Phone},
			"Meta":       &types.AttributeValueMemberM{Value: meta},
		},
	}

	_, err := repo.client.PutItem(context.TODO(), input)
	if err != nil {
		return errors.Wrap(err, "failed to put item")
	}

	return nil
}
