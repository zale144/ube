package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

// Pipeline is a wrapper over an UBE multi event pipeline
type Pipeline struct {
	entity    model.Entity
	actions   []action
	afterEach []action
}

// EventProcessingResult - comment placeholder
type EventProcessingResult struct {
	Status         string                 `json:"status"`
	Errors         []string               `json:"errors,omitempty"`
	BusinessEvents []model.PipelineMedium `json:"-"`
}

const (
	StatusSucceeded       = "Succeeded"
	StatusFailed          = "Failed"
	StatusPartiallyFailed = "PartiallyFailed"
)

// Option is a func type abstraction of a pipeline action
type Option func(*Pipeline)

// NewPipeline uses functional options to construct a Pipeline: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
func NewPipeline(entity model.Entity, options ...Option) *Pipeline {
	if entity == nil {
		zap.L().Fatal("entity is not provided")
	}

	p := &Pipeline{
		entity:  entity,
		actions: []action{},
	}

	for _, option := range options {
		option(p)
	}

	if len(p.actions) == 0 {
		zap.L().Fatal("no actions were provided for the pipeline")
	}

	return p
}

// each action batch?
// file to business events
// business events to file publishing

// InvokePipeline invokes a pipeline TODO: test
func (p *Pipeline) InvokePipeline(ctx context.Context, inputs ...model.Input) (EventProcessingResult, error) {
	if len(inputs) == 0 {
		return EventProcessingResult{}, errors.New("nice try passing empty inputs")
	}

	bes, err := model.InputsToBusinessEvents(inputs, p.entity)
	if err != nil {
		return EventProcessingResult{}, fmt.Errorf("perhaps you want to take another look at your inputs: %w", err)
	}

	var (
		result     EventProcessingResult
		finalError error
	)

	// execute the pipeline actions sequentially
	for idx, act := range p.actions {
		zap.L().Info("Starting pipeline action for all business events. If everything goes well, there will be cake.",
			zap.String("action", act.Name()), zap.Int("action batchsize", act.BatchSize()))

		if act.IsAsync() {
			processAsync(ctx, bes, act, idx)
		} else {
			processSync(ctx, bes, act, idx)
		}

		zap.L().Info("Finished pipeline action for all business events", zap.String("action", act.Name()))

		for _, postAct := range p.afterEach {
			postAct.Process(ctx, bes...)
		}
	}

	zap.L().Info("At last! Finished processing all business events. Let's see how many of them made it to the end.")

	for _, beI := range bes {
		be, ok := beI.(model.PipelineMedium)
		if !ok {
			zap.L().Error("business event is not of type model.PipelineMedium")
			continue
		}
		result.BusinessEvents = append(result.BusinessEvents, be)
	}

	for _, be := range bes {
		if be.GetError() != nil {
			result.Errors = append(result.Errors, be.GetError().Error())
			finalError = multierror.Append(finalError, fmt.Errorf("another one bites the dust. ID: '%s'; err: %w", be.GetID(), be.GetError()))
		}
	}

	errs, events := len(result.Errors), len(bes)
	if errs == 0 {
		zap.L().Info("Congratulations, you made it! All events have been processed flawlessly. Here's your cake: 🍰")
		result.Status = StatusSucceeded
		return result, finalError
	}

	if errs > 0 && errs < events {
		zap.L().Info("Not great, not terrible. Have fun debugging.", zap.Int("events", events), zap.Int("errors", errs))
		result.Status = StatusPartiallyFailed // 1 or more pipeline actions failed
	}

	if errs == events {
		zap.L().Info("We really admire your client's patience. This time everything failed.")
		result.Status = StatusFailed // all pipeline actions failed
	}

	return result, finalError
}

func processSync(ctx context.Context, bes []model.Medium, action action, actionIdx int) {
	lb := len(bes)

	batchSize := action.BatchSize()
	if batchSize == 0 {
		batchSize = 1
	}

	for i := 0; i < lb; i += batchSize {
		j := i + batchSize
		if j > lb || lb == 1 {
			j = lb
		}

		processBatchAction(ctx, bes[i:j], action, actionIdx)
	}
}

func processAsync(ctx context.Context, bes []model.Medium, action action, actionIdx int) {
	wg := sync.WaitGroup{}
	lb := len(bes)

	batchSize := action.BatchSize()
	if batchSize == 0 {
		batchSize = 1
	}

	for i := 0; i < lb; i += batchSize {
		i := i
		wg.Add(1)

		j := i + batchSize
		if j > lb || lb == 1 {
			j = lb
		}

		go func(beb []model.Medium, idx int) {
			defer wg.Done()
			processBatchAction(ctx, beb, action, idx)
		}(bes[i:j], actionIdx)
	}

	wg.Wait()
}

func processBatchAction(ctx context.Context, bes []model.Medium, action action, actionIdx int) {
	if action == nil {
		zap.L().Error("action is nil")
		return
	}

	var toProcess []model.Medium
	// add processable events
	for _, be := range bes { // some might be retries
		if isEventProcessable(be, action, actionIdx) {
			toProcess = append(toProcess, be)
		}
	}
	// if nothing to process - return
	if len(toProcess) == 0 {
		zap.L().Info("No business events to process in current batch.", zap.String("action", action.Name()))
		return
	}
	// process events
	action.Process(ctx, toProcess...)
	// handle errors
	for _, be := range toProcess {
		handleActionError(be.(model.PipelineMedium), action, actionIdx)
	}
}

func handleActionError(be model.PipelineMedium, action action, actionIdx int) {
	be.SetPreviousActionMandate(action.FailureMandate())
	be.SetPreviousAction(actionIdx) // what if redeploying ?

	if be.GetError() == nil {
		be.SetEventProcessedTime(model.Now())
		// clean up from previous re-publish
		if be.GetRepublishAttempt() != nil {
			be.SetRepublishAttempt(nil)
		}
		return
	}

	errText := ""
	switch action.FailureMandate() {
	case model.LogFailureAndContinue:
		errText = "Sorry things didn't work out for your little pipeline action. But, we're not stopping the show because of this tiny hiccup."
	case model.StopFurtherProcessing:
		errText = "Your pipeline action messed up big time, so no further processing of this event will ever happen. Ever."
	case model.StopAndRaiseError:
		errText = "👏...👏...👏... Everything exploded, so not only we're stopping the pipeline altogether, we're also letting your client know!"
	case model.ProcessOnlyCriticalActions:
		errText = "We are sorry to inform you that we will only continue with the critical actions for this particular event."
	case model.StopAndRetry:
		if be.GetRepublishAttempt() == nil {
			attempt := 1
			be.SetRepublishAttempt(&attempt)
		}

		errText = "Tell you what. We will give you another shot at this."
	}

	zap.L().Error(errText, zap.String("action", action.Name()), zap.Error(be.GetError()))
}

func isEventProcessable(beI model.Medium, action action, actionIdx int) bool {
	be := beI.(model.PipelineMedium)
	if be.GetError() == nil {
		// if is re-publish - check if we are at the previously failed action
		if be.GetRepublishAttempt() != nil {
			if actionIdx != be.GetPreviousAction() {
				return false
			}
		}
		return true
	}

	if be.GetPreviousActionMandate() == model.LogFailureAndContinue {
		zap.L().Error("failed to execute action", zap.Error(be.GetError()))
		// start clean slate for next action
		be.SetError(nil)
		return true
	} else if be.GetPreviousActionMandate() == model.StopAndRaiseError && action.Name() == "EventAlerter" {
		// TODO: alert!
		zap.L().Error("ALERT: failed to execute action", zap.Error(be.GetError()))
	} else if be.GetPreviousActionMandate() == model.ProcessOnlyCriticalActions && action.IsCritical() {
		return true
	}

	return false
}

func (p *Pipeline) GetActions() []action {
	return p.actions
}

func (p *Pipeline) SetActions(actions []action) {
	p.actions = actions
}
