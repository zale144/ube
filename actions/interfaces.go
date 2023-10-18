package actions

import (
	"context"
	"io"

	"github.com/zale144/ube/model"
)

//go:generate mockgen -source=./interfaces.go -package actions -destination=./mocks.go

type IAcker interface {
	AckMessages(ctx context.Context, msg ...model.Input) error
}

type IPublisher interface {
	PublishEvents(ctx context.Context, msg ...model.Input) error
}

type IRepublisher interface {
	IAcker
	IPublisher
}

type IUploader interface {
	UploadFile(ctx context.Context, key string, body io.Reader) error
}

type IDownloader interface {
	DownloadFile(ctx context.Context, key string, body io.Writer) error
	DownloadFileFromBucket(ctx context.Context, bucket, key string, body io.Writer) error
}

type IService interface {
	Execute(ctx context.Context, bes ...model.Medium) error
}

type IRepository interface {
	GetEntity(context.Context, model.Key, interface{}) error
	EntityExists(context.Context, model.Key) (bool, error)
	SaveEntities(ctx context.Context, entity ...model.Entity) error
}
