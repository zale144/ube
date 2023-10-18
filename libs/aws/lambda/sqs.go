package lambda

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/zale144/ube/model" // TODO: decouple
)

type (
	// SQSLambdaFn is the standard AWS SQS lambda handler signature
	SQSLambdaFn func(context.Context, events.SQSEvent) error
	// eventHandler provides a standard event handler signature
	eventHandler func(context.Context, *model.InputEvent) error
)

// SQSLambda wraps the pipeline.EventHandler signature with an AWS SQS Lambda handler signature
func SQSLambda(handle eventHandler) SQSLambdaFn {
	return func(ctx context.Context, event events.SQSEvent) error {
		var messages []model.Input

		for i := range event.Records {
			m := &event.Records[i]

			messages = append(messages, &model.Message{
				ID:        m.MessageId,
				Reference: m.ReceiptHandle,
				Body:      []byte(m.Body),
				SourceURI: m.EventSourceARN,
			})
		}

		innerEvent := model.NewInputEvent(messages)

		return handle(ctx, innerEvent)
	}
}
