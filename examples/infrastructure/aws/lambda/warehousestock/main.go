package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	warehousestock "github.com/zale144/ube/examples/internal/warehouseStock"
	"github.com/zale144/ube/handler"
	"github.com/zale144/ube/libs/aws/dynamodb"
	"github.com/zale144/ube/libs/aws/lambda"
	"github.com/zale144/ube/libs/aws/s3"
	"github.com/zale144/ube/libs/aws/sqs"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	uploader := s3.NewUploaderFromEnv("AWS_REGION", "S3_BUCKET")
	warehouseStockRepo := dynamodb.NewDynamoDBFromEnv("DB_WAREHOUSESTOCK_TABLE_NAME")
	publisher := sqs.NewQueueFromEnv("SQS_BQ_QUEUE_URL")
	queue := sqs.NewQueueFromEnv("SQS_QUEUE_URL")

	whs := &warehousestock.UBEModel{}
	pl := whs.Pipeline(uploader, warehouseStockRepo, publisher, queue)
	h := handler.NewEventHandler(pl, queue)
	awslambda.Start(lambda.SQSLambda(h.Handle))
}
