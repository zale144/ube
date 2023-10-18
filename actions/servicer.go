package actions

import (
	"context"
	"fmt"

	"github.com/zale144/ube/model"
)

// Service is a wrapper for injecting business logic into a pipeline via BaseOption
type Service struct {
	service IService
	Base
}

func (e Service) DepCallNames() []string {
	return []string{"Execute"}
}

func (Service) Name() string {
	return "Service"
}

// Servicer constructs a new action with the Base action
func Servicer(service IService, options ...BaseOption) *Service {
	svc := &Service{
		service: service,
		Base: Base{
			batchSize:      100,
			failureMandate: model.StopFurtherProcessing,
			critical:       true,
		},
	}

	for _, opt := range options {
		opt(&svc.Base)
	}

	return svc
}

// Process implements the action interface in UBE, executes the underlying embedded device
func (e Service) Process(ctx context.Context, bes ...model.Medium) {
	var err error
	if err = e.service.Execute(ctx, bes...); err != nil {
		err = fmt.Errorf("execute business service fail: %w", err)
	}

	for i := range bes {
		bes[i].SetError(err)
	}
}
