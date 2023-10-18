package main

import (
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/zale144/ube/examples/internal/car"
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
	carRepo := dynamodb.NewDynamoDBFromEnv("DB_CAR_TABLE_NAME")
	publisher := sqs.NewQueueFromEnv("SQS_BQ_QUEUE_URL")
	ack := sqs.NewQueueFromEnv("SQS_QUEUE_URL")

	carMdl := &car.UBEModel{}
	pl := carMdl.Pipeline(uploader, carRepo, publisher)
	h := handler.NewEventHandler(pl, ack)
	awslambda.Start(lambda.SQSLambda(h.Handle))
}
