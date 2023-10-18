package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zale144/ube/model"
)

// Republish is a wrapper for injecting a retrier into a pipeline
type Republish struct {
	republisher IRepublisher
	maxAttempts int
	Base
}

// Republisher constructs a new Republish
func Republisher(republisher IRepublisher, maxAttempts int, options ...BaseOption) *Republish {
	r := &Republish{
		republisher: republisher,
		Base: Base{
			batchSize:      100,
			failureMandate: model.StopFurtherProcessing,
		},
		maxAttempts: maxAttempts,
	}

	for _, opt := range options {
		opt(&r.Base)
	}

	return r
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (r *Republish) Process(ctx context.Context, bes ...model.Medium) {
	for i, iBe := range bes {
		be, ok := iBe.(model.PipelineMedium)
		if !ok {
			iBe.SetError(fmt.Errorf("expected PipelineMedium, got %T", iBe))
			continue
		}
		if be.GetPreviousActionMandate() == model.StopAndRetry {
			if be.GetRepublishAttempt() == nil || (be.GetRepublishAttempt() != nil && *be.GetRepublishAttempt() >= r.maxAttempts) {
				// TODO: off to the dead letter queue!
				// you had your chance
				continue
			}

			if err := r.republish(ctx, be); err != nil {
				bes[i].SetError(fmt.Errorf("execute business service fail: %w", err))
			}
		}
	}
}

func (r *Republish) republish(ctx context.Context, be model.PipelineMedium) error {
	if r.republisher == nil {
		return fmt.Errorf("re-publisher is not set for the pipeline")
	}

	if be.GetRepublishAttempt() != nil {
		be.IncrementRepublishAttempt()
	}

	be.SetPreviousActionMandate(0)

	msg, err := toMessage(be, false)
	if err != nil {
		return fmt.Errorf("convert business event to message '%s' fail: %w", be.GetID(), err)
	}

	m := make(map[string]interface{})

	if err = json.Unmarshal(msg.Body, &m); err != nil {
		return fmt.Errorf("convert business event to map '%s' fail: %w", be.GetID(), err)
	}

	m["previous_action"] = be.GetPreviousAction()

	jsn, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("convert map to message '%s' fail: %w", be.GetID(), err)
	}

	msg.Body = jsn

	// TODO: find a way to wait until picking up this message
	if err = r.republisher.PublishEvents(ctx, msg); err != nil {
		return fmt.Errorf("re-publish business event '%s' fail: %w", be.GetID(), err)
	}

	ackMsg := &model.Message{
		ID:        be.GetEventID(),
		Reference: be.GetEventReference(),
	}

	// what if original is pointer to file ???
	if err = r.republisher.AckMessages(ctx, ackMsg); err != nil {
		return fmt.Errorf("acknowledge business event '%s' fail: %w", be.GetID(), err)
	}

	be.SetEventID("")
	be.SetEventReference("")
	be.SetPreviousActionMandate(model.StopFurtherProcessing)

	return nil
}

func (*Republish) Name() string {
	return "Republisher"
}

func (r Republish) DepCallNames() []string {
	return []string{"PublishEvents", "AckMessages"}
}
