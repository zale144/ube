package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/zale144/ube/examples/internal/product"
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

	uploader := s3.NewUploaderFromEnv("AWS_REGION", "S3_UPLOAD_BUCKET")
	downloader := s3.NewDownloaderFromEnv("AWS_REGION", "S3_DOWNLOAD_BUCKET")
	prodRepo := dynamodb.NewDynamoDBFromEnv("DB_PRODUCT_TABLE_NAME")
	storeRepo := dynamodb.NewDynamoDBFromEnv("DB_STORE_TABLE_NAME")
	publisher := sqs.NewQueueFromEnv("SQS_BQ_QUEUE_URL")
	queue := sqs.NewQueueFromEnv("SQS_QUEUE_URL")

	prod := &product.Product{}
	pl := prod.Pipeline(uploader, downloader, prodRepo, storeRepo, publisher, queue)
	h := handler.NewEventHandler(pl, queue)
	awslambda.Start(lambda.SQSLambda(h.Handle))
}
