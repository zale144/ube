package sqs

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/model"
)

// publisher is an interface that describes the ability to publish
// an ens.Payload to a message queue/stream.
type publisher interface {
	PublishEvents(ctx context.Context, events ...model.Input) error
}

func TestCreatingNewPublisher(t *testing.T) {
	publisher := NewQueue("some_queue")
	assert.NotNil(t, publisher.client, "should load a service for the publisher")
	assert.IsType(t, (*awssqs.SQS)(nil), publisher.client, "should have a service that is an *sqs.SQS")
}

func TestPublisherIsPublisherCompliant(t *testing.T) {
	assert.Implements(
		t,
		(*publisher)(nil),
		new(Queue),
		"Queue should implement publisher",
	)
}

func TestPublish(t *testing.T) {
	publishFixtures := []struct {
		name            string
		publishOutput   *awssqs.SendMessageOutput
		publishError    error
		queueURL        string
		events          []model.Input
		expectedErr     error
		mockExpectation func(*Mockclient)
	}{
		{
			name:        "a successful publish",
			queueURL:    "queuebar",
			events:      []model.Input{&model.Message{ID: "1", Body: []byte("some body")}},
			expectedErr: nil,
			mockExpectation: func(m *Mockclient) {
				m.EXPECT().SendMessageBatchWithContext(
					gomock.Any(),
					&awssqs.SendMessageBatchInput{
						Entries: []*awssqs.SendMessageBatchRequestEntry{
							{
								Id:          aws.String("1"),
								MessageBody: aws.String("some body"),
							},
						},
						QueueUrl: aws.String("queuebar"),
					},
				).Return(&awssqs.SendMessageBatchOutput{}, nil)
			},
		}, {
			name:        "two messages in one event",
			queueURL:    "queuebar",
			events:      []model.Input{&model.Message{ID: "3", Body: []byte("other body")}, &model.Message{ID: "4", Body: []byte("another body")}},
			expectedErr: nil,
			mockExpectation: func(m *Mockclient) {
				m.EXPECT().SendMessageBatchWithContext(
					gomock.Any(),
					&awssqs.SendMessageBatchInput{
						Entries: []*awssqs.SendMessageBatchRequestEntry{
							{
								Id:          aws.String("3"),
								MessageBody: aws.String("other body"),
							}, {
								Id:          aws.String("4"),
								MessageBody: aws.String("another body"),
							},
						},
						QueueUrl: aws.String("queuebar"),
					},
				).Return(&awssqs.SendMessageBatchOutput{}, nil)
			},
		}, {
			name:          "a failed publish",
			publishOutput: nil,
			publishError:  fmt.Errorf("publish broke"),
			queueURL:      "queuebar",
			events:        []model.Input{&model.Message{ID: "2", Body: []byte("what a body")}},
			expectedErr:   errors.New("send messages fail: publish broke"),
			mockExpectation: func(m *Mockclient) {
				m.EXPECT().SendMessageBatchWithContext(
					gomock.Any(),
					&awssqs.SendMessageBatchInput{
						QueueUrl: aws.String("queuebar"),
						Entries: []*awssqs.SendMessageBatchRequestEntry{
							{
								Id:          aws.String("2"),
								MessageBody: aws.String("what a body"),
							},
						},
					},
				).Return(nil, fmt.Errorf("publish broke"))
			},
		},
	}

	for _, test := range publishFixtures {
		test := test
		t.Run(test.name, func(t *testing.T) {
			asserts := assert.New(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mock := NewMockclient(mockCtrl)

			test.mockExpectation(mock)

			publisher := NewQueueWithService("", mock)
			publisher.queueURL = test.queueURL
			var err error
			if len(test.events) == 1 {
				err = publisher.PublishEvents(context.Background(), test.events[0])
			} else {
				err = publisher.PublishEvents(context.Background(), test.events...)
			}

			if test.expectedErr != nil {
				asserts.EqualError(err, test.expectedErr.Error())
			} else {
				asserts.NoError(err)
			}
		})
	}
}

func TestAckMessages(t *testing.T) {
	deletedMessages := []struct {
		Name            string
		QueueURL        string
		ExpectedError   error
		mockExpectation func(*Mockclient)
		input           model.Input
	}{
		{
			Name:          "Deleted message",
			QueueURL:      "queueURL",
			input:         &model.Message{ID: "1", Reference: "handle1"},
			ExpectedError: nil,
			mockExpectation: func(m *Mockclient) {
				m.EXPECT().DeleteMessageBatchWithContext(
					gomock.Any(),
					&awssqs.DeleteMessageBatchInput{
						QueueUrl: aws.String("queueURL"),
						Entries: []*awssqs.DeleteMessageBatchRequestEntry{
							{
								Id:            aws.String("1"),
								ReceiptHandle: aws.String("handle1"),
							},
						},
					}).Return(&awssqs.DeleteMessageBatchOutput{}, nil)
			},
		},
	}
	for _, test := range deletedMessages {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			asserts := assert.New(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mock := NewMockclient(mockCtrl)

			test.mockExpectation(mock)

			queue := NewQueueWithService("", mock)
			queue.queueURL = test.QueueURL

			err := queue.AckMessages(context.Background(), test.input)
			asserts.Equal(test.ExpectedError, err)
		})
	}
}
