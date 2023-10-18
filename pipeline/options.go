package pipeline

import (
	"github.com/zale144/ube/actions"
)

func Action(act action) Option {
	return func(p *Pipeline) {
		p.SetActions(append(p.GetActions(), act))
	}
}

// AfterEach constructs an action that runs after each pipeline action in the list
func AfterEach(action ...action) Option {
	return func(p *Pipeline) {
		p.afterEach = append(p.afterEach, action...)
	}
}

// Enricher constructs a new action with the Enrich action
func Enricher(enrichers ...actions.EnricherOption) Option {
	return Action(actions.Enricher(enrichers...))
}

// InputTransformer constructs a new action with the InputTransform action
func InputTransformer(options ...actions.TransformOption) Option {
	return Action(actions.InputTransformer(options...))
}

// Servicer constructs a new action with the Servicer action
func Servicer(service actions.IService) Option {
	return Action(actions.Servicer(service))
}

// Uploader constructs a new action with the Uploader action
func Uploader(uploader actions.IUploader, options ...actions.BaseOption) Option {
	return Action(actions.Uploader(uploader, options...))
}

// Persister constructs a new action with the Persist action
func Persister(repo actions.IRepository, options ...actions.BaseOption) Option {
	return Action(actions.Persister(repo, options...))
}

// Publisher constructs a new action with the Publisher action
func Publisher(publisher actions.IPublisher, options ...actions.BaseOption) Option {
	return Action(actions.Publisher(publisher, options...))
}

// Republisher constructs a new action with the Republisher action
func Republisher(republisher actions.IRepublisher, maxAttempts int, options ...actions.BaseOption) Option {
	return AfterEach(actions.Republisher(republisher, maxAttempts, options...))
}
