package dv_dynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/timhugh/digitalvenue/util/core"
	"strings"
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

func (repo *Repository) GetTickets(tenantID string, orderID string) ([]*core.Ticket, error) {
	output, err := repo.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(repo.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: PrefixID("Tenant", tenantID)},
			":sk": &types.AttributeValueMemberS{Value: PrefixID("Order", orderID) + "#Ticket#"},
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get Tickets batch")
	}

	var ticketDTOs []ticketDTO
	err = attributevalue.UnmarshalListOfMaps(output.Items, &ticketDTOs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal Ticket attributes")
	}

	tickets := make([]*core.Ticket, len(ticketDTOs))
	for i, ticketDTO := range ticketDTOs {
		skParams := strings.Split(ticketDTO.SK, "#")
		if len(skParams) != 4 {
			return nil, errors.New("invalid ticket SK format")
		}
		tickets[i] = &core.Ticket{
			ID:         skParams[3],
			TenantID:   tenantID,
			OrderID:    orderID,
			CustomerID: ticketDTO.CustomerID,
			Name:       ticketDTO.Name,
			QRCodeURL:  ticketDTO.QRCodeURL,
		}
	}

	return tickets, nil
}
