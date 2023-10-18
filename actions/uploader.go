package actions

import (
	"bytes"
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

// Upload is a wrapper for uploading files in a pipeline
type Upload struct {
	uploader IUploader
	Base
}

// Uploader constructs a new Uploader
func Uploader(uploader IUploader, options ...BaseOption) *Upload {
	upl := &Upload{
		uploader: uploader,
		Base: Base{
			batchSize:      100,
			failureMandate: model.LogFailureAndContinue,
		},
	}

	for _, opt := range options {
		opt(&upl.Base)
	}

	return upl
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e Upload) Process(ctx context.Context, bes ...model.Medium) {
	var counter int

	for i := range bes {
		be := bes[i]
		if _, ok := e.skips[be.GetEventName()]; ok {
			continue
		}

		data := be.GetRawData()
		for j, d := range data {
			file := bytes.NewBuffer(d)
			// TODO: might want to use something like EntityKey for the file
			key := fmt.Sprintf("%s_%d", be.GetID(), j)

			if err := e.uploader.UploadFile(ctx, key, file); err != nil {
				bes[i].SetError(fmt.Errorf("upload business event raw data fail: %w", err))
				break
			}

			zap.L().Info("upload", zap.String("file", key))

			counter++
		}
	}

	zap.L().Info("files uploaded", zap.Int("size", counter))
}

func (e Upload) DepCallNames() []string {
	return []string{"UploadFile"}
}

func (Upload) Name() string {
	return "Uploader"
}
