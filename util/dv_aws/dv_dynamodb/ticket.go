package dv_dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
)

func (repo *Repository) PutTickets(tickets []*core.Ticket) error {
	writes := make([]types.WriteRequest, len(tickets))

	for i, ticket := range tickets {
		pk := "Tenant#" + ticket.TenantID
		sk := fmt.Sprintf("Order#%s#Ticket#%s", ticket.OrderID, ticket.ID)
		writes[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"PK":         &types.AttributeValueMemberS{Value: pk},
					"SK":         &types.AttributeValueMemberS{Value: sk},
					"CustomerID": &types.AttributeValueMemberS{Value: ticket.CustomerID},
					"Name":       &types.AttributeValueMemberS{Value: ticket.Name},
					"QRCodeURL":  &types.AttributeValueMemberS{Value: ticket.QRCodeURL},
				},
			},
		}
	}

	_, err := repo.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			repo.tableName: writes,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to write tickets batch")
	}

	return nil
}
