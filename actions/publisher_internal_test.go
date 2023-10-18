package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

func TestPublisher_Wrong_EmptyEvent(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{}
	action := Publish{}
	action.Process(ctx, be)

	assert.Contains(t, be.Error.Error(), "publish can't handle empty business event")
}

func TestPublisher_Wrong_NoData(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{Event: &model.Event{}}
	action := Publish{}
	action.Process(ctx, be)

	assert.Contains(t, be.Error.Error(), "publish can't handle empty business event")
}

// mockPublisherMarshalIndentWhichReturnsError is used to mock the json.MarshalIndent function
func mockPublisherMarshalIndentWhichReturnsError(interface{}, string, string) ([]byte, error) {
	return []byte(""), errors.New("fake error")
}

func TestPublisher_Wrong_WrongEncoding(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	MarshalIndent = mockPublisherMarshalIndentWhichReturnsError
	defer func() {
		MarshalIndent = json.MarshalIndent
	}()

	action := Publish{}
	action.Process(ctx, be)

	assert.Error(t, be.Error)
	assert.Equal(t, `convert business event BE-12345 to message fail: marshal business event fail: fake error`, be.Error.Error())
}

func TestPublisher_Wrong_FailToPublish(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := NewMockIPublisher(ctrl)

	gomock.InOrder(
		mockPublisher.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).Return(fmt.Errorf("could not connect to broker")),
	)

	action := Publish{publisher: mockPublisher}
	action.Process(ctx, be)

	assert.Error(t, be.Error)
	assert.Equal(t, `publish message fail: could not connect to broker`, be.Error.Error())
}

func TestPublisher_Good(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := NewMockIPublisher(ctrl)

	gomock.InOrder(
		mockPublisher.EXPECT().PublishEvents(gomock.Any(), gomock.Any()).Return(nil),
	)

	action := Publish{publisher: mockPublisher}
	action.Process(ctx, be)

	assert.NoError(t, be.Error)
}
