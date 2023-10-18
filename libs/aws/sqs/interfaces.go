package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//go:generate mockgen --source=interfaces.go -package sqs -destination=mocks.go

type client interface {
	SendMessageBatchWithContext(
		aws.Context,
		*sqs.SendMessageBatchInput,
		...request.Option,
	) (*sqs.SendMessageBatchOutput, error)
	CreateQueueWithContext(
		aws.Context,
		*sqs.CreateQueueInput,
		...request.Option,
	) (*sqs.CreateQueueOutput, error)
	DeleteMessageBatchWithContext(
		aws.Context,
		*sqs.DeleteMessageBatchInput,
		...request.Option,
	) (*sqs.DeleteMessageBatchOutput, error)
	ReceiveMessageWithContext(
		aws.Context,
		*sqs.ReceiveMessageInput,
		...request.Option,
	) (*sqs.ReceiveMessageOutput, error)
	PurgeQueueWithContext(
		aws.Context,
		*sqs.PurgeQueueInput,
		...request.Option,
	) (*sqs.PurgeQueueOutput, error)
}
