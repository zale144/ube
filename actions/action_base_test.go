package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zale144/ube/model"
)

func TestSkip(t *testing.T) {
	skips := []string{"CreateProduct"}
	want := &Base{
		skips: map[string]struct{}{"CreateProduct": {}},
	}
	base := &Base{}
	Skip(skips...)(base)
	assert.Equalf(t, want, base, "Skip()")
}

func TestFailureMandate(t *testing.T) {
	mandate := model.LogFailureAndContinue
	want := &Base{
		failureMandate: mandate,
	}
	base := &Base{}
	FailureMandate(mandate)(base)
	assert.Equalf(t, want.FailureMandate(), base.FailureMandate(), "FailureMandate()")
}

func TestCritical(t *testing.T) {
	want := &Base{
		critical: true,
	}
	base := &Base{}
	Critical()(base)
	assert.Equalf(t, want.IsCritical(), base.IsCritical(), "Critical()")
}

func TestBatchSize(t *testing.T) {
	batchSize := 13
	want := &Base{
		batchSize: 13,
	}
	base := &Base{}
	BatchSize(batchSize)(base)
	assert.Equalf(t, want.BatchSize(), base.BatchSize(), "BatchSize()")
}

func TestAsync(t *testing.T) {
	want := &Base{
		async: true,
	}
	base := &Base{}
	Async()(base)
	assert.Equalf(t, want.IsAsync(), base.IsAsync(), "Async()")
}
