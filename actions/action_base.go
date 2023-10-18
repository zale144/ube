package actions

import "github.com/zale144/ube/model"

type Base struct {
	critical       bool
	failureMandate model.ActionMandate
	async          bool
	batchSize      int
	skips          map[string]struct{}
}

type BaseOption func(p *Base)

// Skip sets a variadic list of parameter skip event names
func Skip(skips ...string) BaseOption {
	return func(a *Base) {
		skipMap := make(map[string]struct{})
		for _, s := range skips {
			skipMap[s] = struct{}{}
		}
		a.skips = skipMap
	}
}

func FailureMandate(mandate model.ActionMandate) BaseOption {
	return func(a *Base) {
		a.failureMandate = mandate
	}
}

func Critical() BaseOption {
	return func(a *Base) {
		a.critical = true
	}
}

func BatchSize(batchSize int) BaseOption {
	return func(a *Base) {
		a.batchSize = batchSize
	}
}

func Async() BaseOption {
	return func(a *Base) {
		a.async = true
	}
}

// IsCritical marks if an action is critical
func (a Base) IsCritical() bool {
	return a.critical
}

func (a Base) BatchSize() int {
	return a.batchSize
}

func (a Base) IsAsync() bool {
	return a.async
}

// FailureMandate how to handle failure
func (a Base) FailureMandate() model.ActionMandate {
	return a.failureMandate
}
