package pipeline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/zale144/ube/actions"
	"github.com/zale144/ube/model"
)

type product struct {
	productKey
	AnotherOne int
}

type productKey struct {
	SomeField string
}

func (p productKey) PK() string {
	return p.SomeField
}

func (p product) GetKey() model.Key {
	return p.productKey
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestPipeline_FailedTransformer(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	bes := make([]model.Medium, 0)
	be := &model.BusinessEvent{
		Event:    &model.Event{},
		Entities: []model.Entity{&product{AnotherOne: 1}},
	}

	bes = append(bes, be)

	processBatchAction(ctx, bes, actions.InputTransformer(), 0)
	require.Error(t, be.Error)
	assert.Equal(t, `failed to transform business events: unmarshal input body fail: unexpected end of JSON input`, be.Error.Error())
}

func TestPipeline_nil_action(t *testing.T) {
	logger := zap.NewExample()
	defer func() { _ = logger.Sync() }()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	ctx := context.Background()
	bes := make([]model.Medium, 0)
	be := &model.BusinessEvent{}
	bes = append(bes, be)

	processBatchAction(ctx, bes, nil, 0)
	assert.NoError(t, be.Error)
	println()
}
