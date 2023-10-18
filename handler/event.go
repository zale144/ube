package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
	pl "github.com/zale144/ube/pipeline"
)

// EventHandler is the implementation of an event handler
type EventHandler struct {
	pipeline iPipeline
	acker    actions.IAcker
	result   pl.EventProcessingResult
}

// NewEventHandler creates an event handler built with the injected dependencies
func NewEventHandler(pipeline iPipeline, acker actions.IAcker) EventHandler {
	return EventHandler{
		pipeline: pipeline,
		acker:    acker,
	}
}

// Handle handles an event by iterating over its messages
func (p *EventHandler) Handle(ctx context.Context, ev *model.InputEvent) error {
	var (
		ackMsgs   []model.Input
		resultErr error
	)

	resPre, err := p.pipeline.InvokePipeline(ctx, ev.Inputs()...)
	if err != nil {
		zap.L().Error("error while processing record", zap.Error(err))
	}

	p.result = resPre

	for _, e := range resPre.Errors {
		// TODO: make test for this
		resultErr = multierror.Append(resultErr, errors.New(e))
		zap.L().Error("error while processing record", zap.String("error", e))
	}

	// TODO: do we always ack messages, or just when we have a retry mechanism, or when there is no point?
	for _, be := range resPre.BusinessEvents {
		if be.GetEventID() == "" || be.GetEventReference() == "" { // maybe came from pointer to file message
			continue
		}
		msg := &model.Message{
			ID:        be.GetEventID(),
			Reference: be.GetEventReference(),
		}
		ackMsgs = append(ackMsgs, msg)
	}

	if len(ackMsgs) > 0 {
		if err = p.acker.AckMessages(ctx, ackMsgs...); err != nil {
			return fmt.Errorf("acknowledge message fail: %w", err)
		}
	}

	return resultErr
}

func (p EventHandler) GetResult() pl.EventProcessingResult {
	return p.result
}
