package actions

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

type entityObj struct {
	ID   string
	Text string
}

func (eo entityObj) PK() string {
	return fmt.Sprintf("%s", eo.ID)
}

// GetKey returns the custom composed key to the entity
func (eo entityObj) GetKey() model.Key {
	return eo
}

func TestPersister_Wrong_EmptyEvent(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{}
	action := Persist{}
	action.Process(ctx, be)

	assert.Contains(t, be.Error.Error(), "persist can't handle empty business event")
}

func TestPersister_Wrong_NoData(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{Event: &model.Event{}}
	action := Persist{}
	action.Process(ctx, be)

	assert.Contains(t, be.Error.Error(), "persist can't handle empty business event")
}

func TestPersister_Wrong_FailToPersist(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	entity := entityObj{Text: "Zale144"}

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, Entities: []model.Entity{entity}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockIRepository(ctrl)

	gomock.InOrder(
		mockRepository.EXPECT().SaveEntities(gomock.Any(), gomock.Any()).Return(fmt.Errorf("could not connect to database")),
	)

	action := Persist{repo: mockRepository}
	action.Process(ctx, be)

	assert.Error(t, be.Error)
	assert.Equal(t, `persist business event fail: could not connect to database`, be.Error.Error())
}

func TestPersister_Good(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	entity := entityObj{Text: "Zale144"}

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, Entities: []model.Entity{entity}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := NewMockIRepository(ctrl)

	gomock.InOrder(
		mockRepository.EXPECT().SaveEntities(gomock.Any(), gomock.Any()).Return(nil),
	)

	action := Persist{repo: mockRepository}
	action.Process(ctx, be)

	assert.NoError(t, be.Error)
}
