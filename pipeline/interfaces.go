package pipeline

import (
	"context"

	"github.com/zale144/ube/model"
)

// action is an interface for any action that you might want to
// take related to a specific event. They are strung together to
// create Pipelines
type action interface {
	Name() string
	// Process performs the logic for the pipeline action
	Process(ctx context.Context, bes ...model.Medium)
	// IsCritical is a flag that determines whether the action type is considered critical
	// during pipeline processing
	IsCritical() bool
	// FailureMandate is an instruction to the pipeline processor of what to do if this action fails
	FailureMandate() model.ActionMandate

	DepCallNames() []string

	BatchSize() int
	IsAsync() bool
}
