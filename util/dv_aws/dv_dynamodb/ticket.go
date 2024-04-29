package dv_dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
)

type ticketDTO struct {
	PK         string
	SK         string
	CustomerID string
	Name       string
	QRCodeURL  string
}

func (repo *Repository) PutTickets(tickets []*core.Ticket) error {
	writes := make([]types.WriteRequest, len(tickets))

	for i, ticket := range tickets {
		dto := &ticketDTO{
			PK:         PrefixID("Tenant", ticket.TenantID),
			SK:         fmt.Sprintf("Order#%s#Ticket#%s", ticket.OrderID, ticket.ID),
			CustomerID: ticket.CustomerID,
			Name:       ticket.Name,
			QRCodeURL:  ticket.QRCodeURL,
		}
		item, err := attributevalue.MarshalMap(dto)
		if err != nil {
			return errors.Wrap(err, "failed to marshal Ticket attributes")
		}
		item["Type"] = &types.AttributeValueMemberS{Value: "Ticket"}

		writes[i] = types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
	}

	_, err := repo.client.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			repo.tableName: writes,
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to batch write Tickets")
	}

	return nil
}
