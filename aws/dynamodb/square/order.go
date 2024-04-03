package square

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/timhugh/digitalvenue/square"
	"os"
)

type OrderRepositoryConfig struct {
	TableName string
}

func NewOrderRepositoryConfig() OrderRepositoryConfig {
	return OrderRepositoryConfig{
		TableName: os.Getenv(SquareOrdersTableName),
	}
}

type OrderRepository struct {
	tableName string
	client    *dynamodb.Client
}

func NewOrderRepository(config OrderRepositoryConfig, client *dynamodb.Client) square.OrdersRepository {
	return OrderRepository{
		tableName: config.TableName,
		client:    client,
	}
}

func (repo OrderRepository) Create(order square.Order) error {
	putItemInput := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			SquareOrderID:    &types.AttributeValueMemberS{Value: order.SquareOrderID},
			SquareCustomerID: &types.AttributeValueMemberS{Value: order.SquareCustomerID},
		},
		TableName: aws.String(repo.tableName),
	}

	_, err := repo.client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return err
	}

	return nil
}
