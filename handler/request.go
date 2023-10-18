package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

type (
	// APIOUTTransform defines a transform function for converting an event processing result to an API response
	APIOUTTransform func(result *pl.EventProcessingResult) (*model.Response, error)
)

// RequestHandler is the implementation of a request handler
type RequestHandler struct {
	pipeline        iPipeline
	outputTransform APIOUTTransform
}

// NewRequestHandler creates a request handler built with the injected dependencies
func NewRequestHandler(pipeline iPipeline, outputTransform APIOUTTransform) RequestHandler {
	return RequestHandler{
		pipeline:        pipeline,
		outputTransform: outputTransform,
	}
}

// Handle handles a stateless request
func (p RequestHandler) Handle(ctx context.Context, req *model.Request) (*model.Response, error) {
	resPre, err := p.pipeline.InvokePipeline(ctx, req)
	if err != nil {
		zap.L().Error("error while processing record", zap.Error(err))
		err = fmt.Errorf("process record fail: %w", err)
		return &model.Response{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-type": "application/json"},
			Body:       fmt.Errorf(`{"message": "%s"}`, err).Error(),
		}, nil
	}

	var resultErr error
	for _, e := range resPre.Errors {
		resultErr = multierror.Append(resultErr, errors.New(e))
		zap.L().Error("error while processing record",
			zap.String("id", req.GetID()),
			zap.String("error", e))
	}

	resp, err := p.outputTransform(&resPre)
	if err != nil {
		return resp, fmt.Errorf("transform pipeline output to response fail: %w", err)
	}

	return resp, resultErr
}
