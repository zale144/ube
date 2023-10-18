package sqs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/zale144/ube/model" // TODO: decouple
)

const (
	defaultTimeoutSec          = 15
	defaultMaxNumberOfMessages = 10
	defaultWaitTimeSec         = 2
)

// Queue allows for the publishing of ens.Payload event onto an SQS
// queue. This performs message grouping based on the payload.ClientID,
// allowing for fifo to take place across clients, if so desired.
type Queue struct {
	client   client
	name     string
	queueURL string
}

// NewQueue creates a Queue service with the provided queueURL
func NewQueue(url string) *Queue {
	svc := sqs.New(session.Must(session.NewSession()), aws.NewConfig())
	return NewQueueWithService(url, svc)
}

// NewQueueFromEnv creates a Queue service with the provided queueURL
func NewQueueFromEnv(env string) *Queue {
	url := os.Getenv(env)
	if url == "" {
		log.Fatalf("SQS: environment variable '%s' is not set", env)
	}

	svc := sqs.New(session.Must(session.NewSession()), aws.NewConfig())

	return NewQueueWithService(url, svc)
}

// NewQueueWithService creates a Queue service with the provided queueURL,
// name and injected client
func NewQueueWithService(url string, client client) *Queue {
	return &Queue{
		queueURL: url,
		client:   client,
	}
}

// NewQueueWithServiceAndQueue creates a Queue service with the provided queueURL,
// name and injected client, and creates a queue
func NewQueueWithServiceAndQueue(name string, client client) (*Queue, error) {
	// create queue programmatically, to make testing easier
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*defaultTimeoutSec)
	defer cancel()

	q := NewQueueWithService("", client)
	q.name = name

	if err := q.createQueue(ctx); err != nil {
		if err.Error() != errQueueAlreadyExists {
			return q, err
		}
	}

	return q, nil
}

const errQueueAlreadyExists = `queue already exists` // TODO ?

// PublishEvents publishes a batch of events to the queue name that was defined
// during the Queue initialization. This method requires the use of a
// context to determine its running duration. If the events could not be
// event, then an error is returned.
func (q *Queue) PublishEvents(ctx context.Context, messages ...model.Input) error {
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(messages))
	for i := range messages {
		entries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(messages[i].GetID()),
			MessageBody: aws.String(messages[i].GetBody()),
		}
	}

	input := sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(q.queueURL),
	}

	result, err := q.client.SendMessageBatchWithContext(ctx, &input)
	if err != nil {
		return fmt.Errorf("send messages fail: %w", err)
	}

	if result == nil {
		return errors.New("no SQS result produced, unexpected error occurred")
	}

	return nil
}

// AckMessages deletes multiple messages from a sqs once they have been successfully processed otherwise returns an error
func (q *Queue) AckMessages(ctx context.Context, messages ...model.Input) error {
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, len(messages))
	for i := range messages {
		entries[i] = &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(messages[i].GetID()),
			ReceiptHandle: aws.String(messages[i].GetReference()),
		}
	}

	deleteInput := sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(q.queueURL),
	}

	_, err := q.client.DeleteMessageBatchWithContext(ctx, &deleteInput)
	if err != nil {
		return err
	}

	return nil
}

// Poll reads messages from a sqs and executes a handler after
func (q *Queue) Poll(ctx context.Context, next func(context.Context, events.SQSEvent) error) error {
	rec, err := q.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueURL),
		MaxNumberOfMessages: aws.Int64(defaultMaxNumberOfMessages),
		WaitTimeSeconds:     aws.Int64(defaultWaitTimeSec),
	})
	if err != nil {
		return fmt.Errorf("read sqs messsage fail: %w", err)
	}

	msgs := make([]events.SQSMessage, len(rec.Messages))
	for i := range rec.Messages {
		msgs[i] = events.SQSMessage{
			MessageId:     *rec.Messages[i].MessageId,
			ReceiptHandle: *rec.Messages[i].ReceiptHandle,
			Body:          *rec.Messages[i].Body,
		}
	}

	if err = next(ctx, events.SQSEvent{Records: msgs}); err != nil {
		return fmt.Errorf("execute callback function fail: %w", err)
	}

	return nil
}

// Purge empties the queue from all messages
func (q *Queue) Purge(ctx context.Context) error {
	_, err := q.client.PurgeQueueWithContext(ctx, &sqs.PurgeQueueInput{
		QueueUrl: aws.String(q.queueURL),
	})
	if err != nil {
		return fmt.Errorf("purge the queue '%s' fail: %w", q.name, err)
	}

	return nil
}

func (q *Queue) createQueue(ctx context.Context) error {
	rec, err := q.client.CreateQueueWithContext(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(q.name),
	})
	if err != nil {
		return fmt.Errorf("create sqs queue fail: %w", err)
	}

	if rec != nil && rec.QueueUrl != nil {
		q.queueURL = *rec.QueueUrl
	}

	return nil
}
