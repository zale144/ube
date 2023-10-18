package lambda_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/handler"
	"github.com/zale144/ube/libs/aws/lambda"
	"github.com/zale144/ube/model"
	"github.com/zale144/ube/pipeline"
)

type testModel struct {
	model.Base
}

func TestSQSLambda_Good(t *testing.T) {
	tmod := &testModel{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ack := actions.NewMockIAcker(ctrl)
	ack.EXPECT().AckMessages(gomock.Any(), gomock.Any()).AnyTimes()

	mockPublisher := actions.NewMockIPublisher(ctrl)
	mockPublisher.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).Return(nil)

	pline := tmod.Pipeline(mockPublisher)

	h := handler.NewEventHandler(pline, ack)

	lfunc := lambda.SQSLambda(h.Handle)

	ctx := context.Background()
	ev := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: "one",
			},
		},
	}

	err := lfunc(ctx, ev)
	assert.NoError(t, err)
}

func (entity *testModel) Pipeline(publisher actions.IPublisher) *pipeline.Pipeline {
	return pipeline.NewPipeline(
		entity,
		pipeline.Publisher(publisher),
	)
}
