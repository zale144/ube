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

func TestUploader_Good_NoData(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{Event: &model.Event{}}
	action := Upload{}
	action.Process(ctx, be)

	assert.NoError(t, be.Error)
}

func TestUploader_Bad_OneFile(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := NewMockIUploader(ctrl)

	gomock.InOrder(
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("someting is weird")),
	)

	action := Upload{uploader: mockUploader}
	action.Process(ctx, be)

	assert.Error(t, be.Error)
	assert.Contains(t, be.Error.Error(), `upload business event raw data fail: someting is weird`)
}

func TestUploader_Bad_OneOfThreeFiles(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}
	be2 := &model.BusinessEvent{ID: "BE-23456", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}
	be3 := &model.BusinessEvent{ID: "BE-34567", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	// be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`), []byte(`improves`), []byte(`a lot`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := NewMockIUploader(ctrl)

	gomock.InOrder(
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("someting is weird")),
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)

	action := Upload{uploader: mockUploader}
	action.Process(ctx, be, be2, be3)

	assert.Error(t, be.Error)
	assert.Contains(t, be.Error.Error(), `upload business event raw data fail: someting is weird`)

	assert.NoError(t, be2.Error)
	assert.NoError(t, be3.Error)
}

func TestUploader_Good_OneFile(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := NewMockIUploader(ctrl)

	gomock.InOrder(
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)

	action := Upload{uploader: mockUploader}
	action.Process(ctx, be)

	assert.NoError(t, be.Error)
}

func TestUploader_Good_ThreeFiles(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	be := &model.BusinessEvent{ID: "BE-12345", Event: &model.Event{}, RawDataEvent: [][]byte{[]byte(`Zale144`), []byte(`improves`), []byte(`a lot`)}}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := NewMockIUploader(ctrl)

	gomock.InOrder(
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
		mockUploader.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil),
	)

	action := Upload{uploader: mockUploader}
	action.Process(ctx, be)

	assert.NoError(t, be.Error)
}
