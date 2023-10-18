package handler

import (
	"context"

	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

//go:generate mockgen -source=./interfaces.go -package handler -destination=./mocks.go -mock_names=iPipeline=MockPipeline

type iPipeline interface {
	InvokePipeline(ctx context.Context, input ...model.Input) (pl.EventProcessingResult, error)
}
