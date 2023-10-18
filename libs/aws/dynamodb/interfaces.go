package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//go:generate mockgen -destination=./mocks.go -package=dynamodb -source=interfaces.go
//-mock_names dynamoDBAPI=mockDynamoDBAPI timeGenerator=mockTimeGenerator

type dynamoDB interface {
	UpdateItemWithContext(
		ctx aws.Context,
		input *dynamodb.UpdateItemInput,
		opts ...request.Option,
	) (*dynamodb.UpdateItemOutput, error)
	// TransactWriteItemsWithContext(
	// 	ctx aws.Context,
	// 	input *dynamodb.TransactWriteItemsInput,
	// 	opts ...request.Option,
	// ) (*dynamodb.TransactWriteItemsOutput, error)
	GetItemWithContext(
		aws.Context,
		*dynamodb.GetItemInput,
		...request.Option,
	) (*dynamodb.GetItemOutput, error)
	CreateTableWithContext(
		aws.Context,
		*dynamodb.CreateTableInput,
		...request.Option,
	) (*dynamodb.CreateTableOutput, error)
}
