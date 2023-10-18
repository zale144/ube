package actions

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/zale144/ube/model"
)

// Persist is a wrapper for persisting the business event be to a provided data store
type Persist struct {
	repo IRepository
	Base
}

// Persister constructs a new Persist
func Persister(repo IRepository, options ...BaseOption) *Persist {
	pers := &Persist{
		repo: repo,
		Base: Base{
			batchSize:      100,
			failureMandate: model.StopAndRaiseError,
			critical:       true,
		},
	}

	for _, opt := range options {
		opt(&pers.Base)
	}

	return pers
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e Persist) Process(ctx context.Context, bes ...model.Medium) {
	var ents []model.Entity

	for i := range bes {
		be := bes[i]
		if be.GetID() == "" {
			bes[i].SetError(errors.New("persist can't handle empty business event"))
			continue
		}

		if _, ok := e.skips[be.GetEventName()]; ok {
			continue
		}

		ents = append(ents, be.GetEntities()...)
	}

	if len(ents) > 0 {
		// upsert
		var err error
		if err = e.repo.SaveEntities(ctx, ents...); err != nil {
			err = fmt.Errorf("persist business event fail: %w", err)
		}
		for i := range bes {
			bes[i].SetError(err)
		}
	}

	zap.L().Info("entities persisted", zap.Int("size", len(ents)))
}

func (e Persist) DepCallNames() []string {
	return []string{"SaveEntities"}
}

func (e Persist) Name() string {
	return "Persister"
}
