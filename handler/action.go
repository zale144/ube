package handler

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/zale144/ube/model"
)

type EventActionHandler struct {
	action action
	toBE   toBE
}

type toBE func(input model.Input) model.Medium

// NewEventActionHandler creates an action event handler built with the injected dependencies
func NewEventActionHandler(action action, toBE toBE) EventActionHandler {
	return EventActionHandler{
		action: action,
		toBE:   toBE,
	}
}

type action interface {
	Name() string
	// Process performs the logic for the pipeline action
	Process(ctx context.Context, bes ...model.Medium)
}

type (
	InternalHandler func(ctx context.Context, ev *model.InputEvent) error
)

// Handle handles an event by iterating over its messages
func (p *EventActionHandler) Handle(ctx context.Context, ev *model.InputEvent) error {
	if p.action == nil {
		return fmt.Errorf("no action defined")
	}

	if p.toBE == nil {
		return fmt.Errorf("no toBE function defined")
	}

	inputs := ev.Inputs()
	bes := make([]model.Medium, len(inputs))
	for _, in := range inputs {
		bes = append(bes, p.toBE(in))
	}

	p.action.Process(ctx, bes...)

	var err error
	for _, be := range bes {
		if be.GetError() != nil {
			err = multierror.Append(err, fmt.Errorf("error while processing record: %w", be.GetError()))
		}
	}

	return err
}
